package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/rodeorm/shortener/internal/api"
	"github.com/rodeorm/shortener/internal/logger"
	"github.com/rodeorm/shortener/internal/repo"
)

// config выполняет первоначальную конфигурацию
func configurate() (*api.Server, error) {
	flag.Parse()

	var (
		serverAddress, baseURL, fileStoragePath, databaseConnectionString, configName, httpsEnabled, trustedSubnet string
		workerCount, batchSize, queueSize, profileType                                                             int
		err                                                                                                        error
	)
	logger.Initialize("info")

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

	ms, fs, ps := repo.GetStorages(fileStoragePath, databaseConnectionString)
	builder := &api.ServerBuilder{}

	if ps != nil {
		server := builder.SetStorages(ps, ps, ps, ps).
			SetDeleter(workerCount, batchSize, queueSize).
			SetConfig(serverAddress, baseURL, fileStoragePath, databaseConnectionString, httpsEnabled, trustedSubnet).
			SetConfigFromFile(configName).
			SetProfileType(profileType).
			Build()

		return &server, nil

	} else if fs != nil {
		server := builder.SetStorages(fs, fs, nil, fs).
			SetDeleter(workerCount, batchSize, queueSize).
			SetConfig(serverAddress, baseURL, fileStoragePath, databaseConnectionString, httpsEnabled, trustedSubnet).
			SetConfigFromFile(configName).
			SetProfileType(profileType).
			Build()

		return &server, nil
	}
	server := builder.SetStorages(ms, ms, nil, ms).
		SetDeleter(workerCount, batchSize, queueSize).
		SetConfig(serverAddress, baseURL, fileStoragePath, databaseConnectionString, httpsEnabled, trustedSubnet).
		SetConfigFromFile(configName).
		SetProfileType(profileType).
		Build()

	return &server, nil
}
