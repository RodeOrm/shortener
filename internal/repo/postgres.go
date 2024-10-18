package repo

import (
	"context"
	"fmt"
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

func (s postgresStorage) Ping() error {
	return s.DB.Ping()
}

/*
InsertUser принимает идентификатор пользователя

Возвращает по идентификатору уже имеющегося в наличии пользователя, если такового нет, то создает нового и возвращает что пользователь не был авторизован по переданному идентификатору
*/
func (s postgresStorage) InsertUser(Key int) (*core.User, error) {

	// ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	// defer cancel()

	ctx := context.TODO()
	//Ищем пользователя
	err := s.preparedStatements["SelectUser"].GetContext(ctx, &Key, Key)

	//При любой ошибке (нет пользователя с таким ИД или передан 0 в Key) получаем нового
	if err != nil {
		err = s.preparedStatements["InsertUser"].GetContext(ctx, &Key, time.Now().Format(time.DateTime))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", "ошибка при InsertUser", err)
		}
		return &core.User{Key: Key, WasUnathorized: true}, nil
	}
	return &core.User{Key: Key, WasUnathorized: false}, nil
}

/*
	InsertShortURL принимает оригинальный URL, генерирует для него ключ, сохраняет соответствие оригинального URL и ключа.

Возвращает соответствующий сокращенный урл, а также признак того, что url сократили ранее
*/
func (s postgresStorage) InsertURL(URL, baseURL string, user *core.User) (*core.URL, error) {
	if !core.CheckURLValidity(URL) {
		return nil, fmt.Errorf("невалидный URL: %s", URL)
	}

	ctx := context.TODO()

	url, err := s.getShortURL(ctx, URL)
	if err != nil {
		return nil, err
	}

	s.preparedStatements["InsertURL"].ExecContext(ctx, url.OriginalURL, url.Key, user.Key)

	return url, nil

}

func (s postgresStorage) getShortURL(ctx context.Context, URL string) (*core.URL, error) {
	url := core.URL{OriginalURL: URL}
	// Смотрим - не сокращали ли урл ранее, если сокращали, то возвращаем ключ для сокращенного
	err := s.preparedStatements["SelectShortURL"].GetContext(ctx, &url.Key, url.OriginalURL)
	if err == nil {
		url.HasBeenShorted = true
		return &url, nil
	}
	// В ином случае получаем новый ключ
	url.Key, err = core.ReturnShortKey(5)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "ошибка при обращении ReturnShortKey из SelectShortURL", err)
	}
	return &url, nil
}

/*
	SelectOriginalURL принимает короткий урл.

Возвращает соответствующий оригинальный урл, признак, что url ранее уже сокращался; признак, что url удален
*/
func (s postgresStorage) SelectOriginalURL(shortURL string) (*core.URL, error) {
	ctx := context.TODO()
	url := core.URL{Key: shortURL}

	err := s.preparedStatements["SelectOriginalURL"].QueryRowContext(ctx, shortURL).Scan(&url.OriginalURL, &url.HasBeenDeleted)

	if err != nil {
		return nil, fmt.Errorf("ошибка в SelectOriginalURL: %v", err)
	}

	url.HasBeenShorted = true

	return &url, nil
}

// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s postgresStorage) SelectUserURLHistory(user *core.User) ([]core.UserURLPair, error) {
	urls := make([]core.UserURLPair, 0, 1)

	err := s.preparedStatements["SelectUserURLHistory"].Select(&urls, user.Key)

	if err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("нет истории для пользователя %d", user.Key)
	}
	return urls, nil
}

func (s postgresStorage) Close() {
	s.DB.Close()
}

func (s postgresStorage) DeleteURLs(URLs []core.URL) error {
	tx := s.DB.MustBegin()
	defer tx.Rollback()

	query := `UPDATE Urls SET isDeleted = true WHERE short = :key AND userID = :user_key`

	for _, update := range URLs {
		_, err := tx.NamedExec(query, update)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("%s: %w: %s: %w", "ошибка при обновлении", err, "ошибка при откате транзакции", rbErr)
			}
			return fmt.Errorf("откат транзакции из-за ошибки при обновлении %s: %d, %w", update.Key, update.UserKey, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", "ошибка при фиксации транзакции", err)
	}
	return nil
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
		return fmt.Errorf("%s: %w", "ошибка при создании таблиц", err)
	}
	return nil
}
