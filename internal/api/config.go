package api

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
)

// Config конфигурация сервера
type Config struct {
	ServerAddress   string `json:"server_address,omitempty"`    // "server_address": "localhost:8080",
	BaseURL         string `json:"base_url,omitempty"`          // "base_url": "http://localhost",
	FileStoragePath string `json:"file_storage_path,omitempty"` // "file_storage_path": "/path/to/file.db",
	DatabaseDSN     string `json:"database_dsn,omitempty"`      //  "database_dsn": "",
	EnableHTTPS     bool   `json:"enable_https,omitempty"`      // "enable_https": true
	IsGivenHTTPS    bool   // Для случаев, когда значение не представлено
}

// Deleter конфигурация сервера для удаления
type Deleter struct {
	WorkerCount int // Количество воркеров, асинхронно удаляющих url
	BatchSize   int // Размер пачки для удаления

	DeleteQueue *Queue //Очередь удаления

}

// ServerBuilder абстракция для создания сервера
type ServerBuilder struct {
	server Server
}

// SetConfigFromFile заполняет конфигурацию данными конфигурационного файла
func (s ServerBuilder) SetConfigFromFile(configName string) ServerBuilder {

	file, err := os.Open(configName)
	if err != nil {
		return s
	}
	defer file.Close()

	var cfg Config
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &cfg)
	}

	if s.server.Config.BaseURL == "" {
		s.server.Config.BaseURL = cfg.BaseURL
	}

	if s.server.Config.DatabaseDSN == "" {
		s.server.Config.DatabaseDSN = cfg.DatabaseDSN
	}

	if s.server.Config.FileStoragePath == "" {
		s.server.Config.FileStoragePath = cfg.FileStoragePath
	}

	if s.server.Config.ServerAddress == "" {
		s.server.Config.ServerAddress = cfg.ServerAddress
	}

	if !s.server.Config.IsGivenHTTPS {
		s.server.Config.EnableHTTPS = cfg.EnableHTTPS
	}

	return s
}

// SetConfig заполняет конфигурацию данными из переменных окружения и флагов
func (s ServerBuilder) SetConfig(sa, bu, fsp, dn, eh string) ServerBuilder {
	s.server.Config = Config{ServerAddress: sa,
		BaseURL:         bu,
		FileStoragePath: fsp,
		DatabaseDSN:     dn,
	}

	enableHTTPS, err := strconv.ParseBool(eh)
	if err != nil {
		s.server.Config.IsGivenHTTPS = false
		return s
	}
	s.server.EnableHTTPS = enableHTTPS
	s.server.IsGivenHTTPS = true
	return s
}

// SetStorages указывает хранилища для сервера
func (s ServerBuilder) SetStorages(url URLStorager, user UserStorager, db DBStorager) ServerBuilder {
	s.server.URLStorage = url
	s.server.UserStorage = user
	s.server.DBStorage = db
	return s
}

// SetDeleter конфигурирует удаление URL
func (s ServerBuilder) SetDeleter(wc, bs, qs int) ServerBuilder {
	s.server.Deleter = Deleter{WorkerCount: wc, BatchSize: bs, DeleteQueue: NewQueue(qs)}
	return s
}

// SetProfileType конфигурирует профилирование
func (s ServerBuilder) SetProfileType(profileType int) ServerBuilder {
	s.server.ProfileType = profileType
	return s
}

// Build возвращает сконфигурированный сервер объект Car.
func (s ServerBuilder) Build() Server {
	return s.server
}
