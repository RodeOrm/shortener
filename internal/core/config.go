package core

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 15 * time.Second
	shutdownTimeout    = 30 * time.Second
)

// Config конфигурация сервера
type Config struct {
	ServerConfig
	DatabaseConfig
	TLSConfig
}

// ServerConfig основные параметры сервера
type ServerConfig struct {
	ServerAddress      string `json:"server_address,omitempty"` // "server_address": "localhost:8080"
	GRPCAddress        string `json:"grpc_address,omitempty"`
	BaseURL            string `json:"base_url,omitempty"`          // "base_url": "http://localhost"
	FileStoragePath    string `json:"file_storage_path,omitempty"` // "file_storage_path": "/path/to/file.db"
	TrustedSubnet      string `json:"trusted_subnet,omitempty"`
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ShutdownTimeout    time.Duration
}

// DatabaseConfig параметры, связанные с СУБД
type DatabaseConfig struct {
	DatabaseDSN string `json:"database_dsn,omitempty"` //  "database_dsn": ""
}

// TLSConfig паарметры, связаные с https
type TLSConfig struct {
	EnableHTTPS  bool `json:"enable_https,omitempty"` // "enable_https": true
	IsGivenHTTPS bool // Для случаев, когда значение не представлено
}

// Deleter конфигурация сервера для удаления
type Deleter struct {
	WorkerCount int // Количество воркеров, асинхронно удаляющих url
	BatchSize   int // Размер пачки для удаления

	DeleteQueue *Queue //Очередь удаления

}

// Configurate выполняет первоначальную конфигурацию
func Configurate(a, b, c, config, d, f, w, s, q, p, bs, t *string) (*Server, error) {

	var (
		serverAddress, baseURL, fileStoragePath, databaseConnectionString, configName, httpsEnabled, trustedSubnet string
		workerCount, batchSize, queueSize, profileType                                                             int
		err                                                                                                        error
	)

	//Адрес запуска HTTP-сервера
	if *a == "" {
		serverAddress = os.Getenv("SERVER_ADDRESS")
		if serverAddress == "" {
			serverAddress = "localhost:8080"
		}
	} else {
		serverAddress = *a
	}

	//Базовый адрес результирующего сокращённого URL
	if *b == "" {
		baseURL = os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
	} else {
		baseURL = *b
	}

	//Имя файла конфигурации должно задаваться через флаг -c/-config или переменную окружения CONFIG
	if *c != "" {
		configName = *c
	} else if *config == "" {
		configName = *config
	} else {
		configName = os.Getenv("CONFIG")
	}

	//Строка подключения к БД
	if *d == "" {
		databaseConnectionString = os.Getenv("DATABASE_DSN")
	} else {
		databaseConnectionString = *d
	}

	fmt.Println(configName)

	//Путь до файла
	if *f == "" {
		fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	} else {
		fileStoragePath = *f
	}

	if *w == "" {
		workerCount = 2
	}

	if *bs == "" {
		batchSize = 3
	}

	if *q == "" {
		queueSize = 10
	}

	if *p == "" {
		profileType = noneProfile
	} else {
		profileType, err = strconv.Atoi(*p)
		if err != nil {
			profileType = noneProfile
		}
	}

	/*
		При передаче флага -s или переменной окружения ENABLE_HTTPS запускайте сервер с помощью метода http.ListenAndServeTLS или tls.Listen.
	*/
	if *s == "" {
		httpsEnabled = os.Getenv("ENABLE_HTTPS")
	} else {
		httpsEnabled = *s
	}

	/**/
	if *t == "" {
		trustedSubnet = os.Getenv("TRUSTED_SUBNET")
	}

	// ms, fs, ps := repo.GetStorages(fileStoragePath, databaseConnectionString)
	builder := &ServerBuilder{}

	server := builder.SetDeleter(workerCount, batchSize, queueSize).
		SetConfig(serverAddress, baseURL, fileStoragePath, databaseConnectionString, httpsEnabled, trustedSubnet).
		SetConfigFromFile(configName).
		SetProfileType(profileType).
		SetTimeOuts(serverReadTimeout, serverWriteTimeout, shutdownTimeout).
		Build()

	return &server, nil
}
