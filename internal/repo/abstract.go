package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/logger"
	"go.uber.org/zap"
)

type AbstractStorage interface {
	/*
		InsertURL принимает оригинальный URL, базовый урл для генерации коротких адресов и пользователя.
		Генерирует уникальный ключ для короткого адреса, сохраняет соответствие оригинального URL и ключа.

		Возвращает соответствующий сокращенный урл, а также признак того, что url сократили ранее
	*/
	InsertURL(URL, baseURL string, user *core.User) (string, bool, error)

	// SelectOriginalURL возвращает оригинальный адрес на основании короткого; признак, что url ранее уже сокращался; признак, что url удален
	SelectOriginalURL(shortURL string) (string, bool, bool, error)

	// InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии, а также значение "отсутствие авторизации по переданному идентификатору"
	InsertUser(Key int) (*core.User, bool, error)

	// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
	SelectUserURLHistory(user *core.User) (*[]core.UserURLPair, error)

	// Массово помечает URL как удаленные. Успешно удалить URL может только пользователь, его создавший.
	DeleteURLs(URL string, user *core.User) (bool, error)

	// Закрыть соединение (только для СУБД)
	CloseConnection()
}

// NewStorage определяет место для хранения данных
func NewStorage(filePath, dbConnectionString string) AbstractStorage {
	var storage AbstractStorage

	storage, err := InitPostgresStorage(dbConnectionString)
	if err == nil {
		return storage
	}

	if filePath != "" {
		storage, err = InitFileStorage(filePath)
		if err == nil {
			return storage
		}
	}
	storage = InitMemoryStorage()
	return storage
}

// InitMemoryStorage создает хранилище данных в оперативной памяти
func InitMemoryStorage() *memoryStorage {
	ots := make(map[string]string)
	sto := make(map[string]string)
	usr := make(map[int]*core.User)
	usrURL := make(map[int]*[]core.UserURLPair)
	storage := memoryStorage{originalToShort: ots, shortToOriginal: sto, users: usr, userURLPairs: usrURL}

	logger.Log.Info("Init storage",
		zap.String("Storage", "Memory storage"),
	)

	return &storage
}

// InitFileStorage создает хранилище данных на файловой системе
func InitFileStorage(filePath string) (*fileStorage, error) {
	usr := make(map[int]*core.User)
	usrURL := make(map[int]*[]core.UserURLPair)
	storage := fileStorage{filePath: filePath, users: usr, userURLPairs: usrURL}
	err := storage.CheckFile(filePath)
	if err != nil {
		return nil, err
	}

	logger.Log.Info("Init storage",
		zap.String("Storage", "File storage"),
	)

	return &storage, nil
}

// InitPostgresStorage создает хранилище данных в БД на экземпляре Postgres
func InitPostgresStorage(connectionString string) (*postgresStorage, error) {
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	delQueue := make(chan string)
	storage := postgresStorage{DB: db, ConnectionString: connectionString, deleteQueue: delQueue, preparedStatements: map[string]*sqlx.Stmt{}}
	err = storage.createTables(ctx)
	if err != nil {
		return nil, err
	}

	err = storage.prepareStatements()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	logger.Log.Info("Init storage",
		zap.String("Storage", "PostgresStorage"),
	)

	return &storage, nil
}

func (s *postgresStorage) prepareStatements() error {

	nstmtSelectUser, err := s.DB.Preparex(`SELECT ID from Users WHERE ID = $1`)
	if err != nil {
		return err
	}

	nstmtInsertUser, err := s.DB.Preparex(`INSERT INTO Users (Name) VALUES ($1) RETURNING ID`)
	if err != nil {
		return err
	}

	nstmtSelectShortURL, err := s.DB.Preparex(`SELECT short from Urls WHERE original = $1`)
	if err != nil {
		return err
	}

	s.preparedStatements["nstmtSelectUser"] = nstmtSelectUser
	s.preparedStatements["nstmtInsertUser"] = nstmtInsertUser
	s.preparedStatements["nstmtSelectShortURL"] = nstmtSelectShortURL
	return nil
}
