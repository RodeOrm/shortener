package repo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rodeorm/shortener/internal/core"
)

type postgresStorage struct {
	DB                 *sqlx.DB    // Драйвер подключения к СУБД
	DBName             string      // Имя БД из конфиг.файла
	ConnectionString   string      // Строка подключения из конфиг.файла
	deleteQueue        chan string // канал для удаления URL
	preparedStatements map[string]*sqlx.Stmt
}

/*
InsertUser принимает идентификатор пользователя

Возвращает по идентификатору уже имеющегося в наличии пользователя, если такового нет, то создает нового и возвращает что пользователь не был авторизован по переданному идентификатору
*/
func (s postgresStorage) InsertUser(Key int) (*core.User, bool, error) {

	ctx := context.TODO()
	var isUnathorized bool

	//Ищем пользователя
	err := s.preparedStatements["nstmtSelectUser"].GetContext(ctx, &Key, Key)

	//При любой ошибке (нет пользователя с таким ИД или передан 0 в Key) получаем нового
	if err != nil {
		isUnathorized = true
		err = s.preparedStatements["nstmtInsertUser"].GetContext(ctx, &Key, time.Now().Format(time.DateTime))
		if err != nil {
			return nil, isUnathorized, err
		}
	}
	return &core.User{Key: Key}, isUnathorized, nil
}

/*
	InsertShortURL принимает оригинальный URL, генерирует для него ключ, сохраняет соответствие оригинального URL и ключа.

Возвращает соответствующий сокращенный урл, а также признак того, что url сократили ранее
*/
func (s postgresStorage) InsertURL(URL, baseURL string, user *core.User) (string, bool, error) {

	if !core.CheckURLValidity(URL) {
		return "", false, fmt.Errorf("невалидный URL: %s", URL)
	}

	ctx := context.TODO()

	var short string

	s.preparedStatements["nstmtSelectShortURL"].GetContext(ctx, &short, URL)
	if short != "" {
		return short, true, nil
	}
	// Вставляем новый URL
	shortKey, err := core.ReturnShortKey(5)
	if err != nil {
		return "", false, err
	}

	arg := map[string]interface{}{
		"original": URL,
		"shortKey": shortKey,
		"userID":   user.Key,
	}

	nstmtInsertURL, args, err := sqlx.Named("INSERT INTO Urls (original, short, userID) SELECT :original, :shortKey, :userID", arg)
	if err != nil {
		return "", false, err
	}

	nstmtInsertURL, args, err = sqlx.In(nstmtInsertURL, args...)
	if err != nil {
		return "", false, err
	}

	nstmtInsertURL = s.DB.Rebind(nstmtInsertURL)
	_, err = s.DB.ExecContext(ctx, nstmtInsertURL, args...)
	if err != nil {
		return "", false, err
	}

	return shortKey, false, nil

}

/*
	SelectOriginalURL принимает короткий урл.

Возвращает соответствующий оригинальный урл, признак, что url ранее уже сокращался; признак, что url удален
*/
func (s postgresStorage) SelectOriginalURL(shortURL string) (string, bool, bool, error) {
	ctx := context.TODO()
	var (
		original    string
		isDeleted   bool
		isShortened bool
	)

	err := s.DB.QueryRowContext(ctx, "SELECT original, isDeleted FROM Urls WHERE short = $1", shortURL).Scan(&original, &isDeleted)
	if err != nil {
		log.Println("SelectOriginalURL", err)
		return "", false, false, err
	}

	isShortened = true

	return original, isShortened, isDeleted, nil
}

// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s postgresStorage) SelectUserURLHistory(user *core.User) (*[]core.UserURLPair, error) {
	urls := make([]core.UserURLPair, 0, 1)
	err := s.DB.Select(&urls, "SELECT original AS origin, short, userID AS userkey FROM Urls WHERE UserID = $1", user.Key)

	if err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("нет истории")
	}
	return &urls, nil
}

func (s postgresStorage) CloseConnection() {
	s.DB.Close()
}

func (s postgresStorage) DeleteURLs(URL string, user *core.User) (bool, error) {
	ch := make(chan string)

	urls := core.GetSliceFromString(URL)

	go func() {
		for _, url := range urls {
			ch <- url
		}
		close(ch)
	}()

	for v := range makeDeletePool(ch) {
		go s.deleteURL(v, user)
	}

	return true, nil
}

func (s postgresStorage) deleteURL(url string, user *core.User) (bool, error) {
	tx, err := s.DB.Begin()

	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(context.TODO(), "UPDATE Urls SET isDeleted = true WHERE short = $1 AND userID = $2")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(context.TODO(), url, user.Key)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func makeDeletePool(inputChs ...chan string) chan string {
	outCh := make(chan string)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan string) {
				defer wg.Done()
				for item := range inputCh {
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func (s postgresStorage) createTables(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS  Users"+
			"("+
			"ID INT GENERATED BY DEFAULT AS IDENTITY"+
			", PRIMARY KEY (ID)"+
			", Name TEXT NULL"+
			")"+
			"; CREATE TABLE IF NOT EXISTS  Urls"+
			"("+
			"ID INT GENERATED BY DEFAULT AS IDENTITY"+
			", PRIMARY KEY (ID)"+
			", isDeleted BOOLEAN NOT NULL DEFAULT False"+
			", UserID	INT  REFERENCES Users (ID) NOT NULL"+
			", Original TEXT NOT NULL "+
			", CorrelationID TEXT NULL"+
			", Short TEXT NOT NULL"+
			");"+
			"CREATE UNIQUE INDEX IF NOT EXISTS url_unique_idx ON Urls (original, UserID) INCLUDE (short);")
	if err != nil {
		log.Println("createTables", err)
		return err
	}
	return nil
}
