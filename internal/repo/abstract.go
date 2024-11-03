// Package repo предназначен для реализации взаимодействия с хранилищами данных
package repo

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/logger"
	"go.uber.org/zap"
)

var (
	ms     *memoryStorage
	fs     *fileStorage
	ps     *postgresStorage
	onceMS sync.Once
	onceFS sync.Once
	oncePS sync.Once
)


// GetStorages определяет реализации для хранения данных
func GetStorages(filePath, dbConnectionString string) (*memoryStorage, *fileStorage, *postgresStorage) {
	logger.Log.Info("Init storage",
		zap.String("Начали процесс выбора хранилища", filePath),
	)
	ps, err := GetPostgresStorage(dbConnectionString)
	if err == nil {
		return nil, nil, ps
	}
	fs, err := GetFileStorage(filePath)
	if err == nil {
		return nil, fs, nil
	}
	logger.Log.Info("Init storage",
		zap.String("хранилище в памяти", ""),
	)
	return GetMemoryStorage(), nil, nil
}

// GetMemoryStorage возвращает хранилище данных в оперативной памяти (создает, если его не было ранее)
func GetMemoryStorage() *memoryStorage {
	onceMS.Do(
		func() {
			ots := make(map[string]string)
			sto := make(map[string]string)
			usr := make(map[int]*core.User)
			usrURL := make(map[int]*[]core.UserURLPair)
			ms = &memoryStorage{originalToShort: ots, shortToOriginal: sto, users: usr, userURLPairs: usrURL}
			logger.Log.Info("Init storage",
				zap.String("Storage", "Memory storage"),
			)
		})
	return ms
}

// GetFileStorage возвращает хранилище данных на файловой системе  (создает, если его не было ранее)
func GetFileStorage(filePath string) (*fileStorage, error) {
	onceFS.Do(
		func() {
			usr := make(map[int]*core.User)
			usrURL := make(map[int]*[]core.UserURLPair)

			fs = &fileStorage{filePath: filePath, users: usr, userURLPairs: usrURL}
			logger.Log.Info("Init storage",
				zap.String("Storage", "File storage"),
			)
		})
	if err := сheckFile(filePath); err != nil {
		logger.Log.Error("can't define file storage",
			zap.Error(err),
		)
		return nil, err
	}
	return fs, nil
}

// GetPostgresStorage возвращает хранилище данных в Postgres (создает, если его не было ранее)
func GetPostgresStorage(connectionString string) (*postgresStorage, error) {
	var (
		dbErr error
		db    *sqlx.DB
	)
	oncePS.Do(
		func() {
			db, dbErr = sqlx.Open("pgx", connectionString)
			if dbErr != nil {
				return
			}
			delQueue := make(chan string)
			ps = &postgresStorage{DB: db, ConnectionString: connectionString, deleteQueue: delQueue, preparedStatements: map[string]*sqlx.Stmt{}}

			ctx := context.TODO()

			if dbErr = ps.createTables(ctx); dbErr != nil {
				return
			}

			if dbErr = ps.createTables(ctx); dbErr != nil {
				return
			}
			if dbErr = ps.prepareStatements(); dbErr != nil {
				return
			}
			logger.Log.Info("Init storage",
				zap.String("Storage", "Postgres storage"),
			)
		})

	if dbErr != nil {
		logger.Log.Error("can't define postgres storage",
			zap.Error(dbErr),
		)
		return nil, dbErr
	}

	return ps, nil
}
