package core

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Server - общий набор атрибутов для http и grpc сервера
type Server struct {
	IdleConnsClosed chan struct{} // Уведомление о завершении работы

	ProfileType int // Тип профилирования (если необходимо)

	URLStorage  URLStorager  // Хранилище данных для URL
	UserStorage UserStorager // Хранилище данных для URL
	DBStorage   DBStorager   // Хранилище данных для DB
	StatStorage StatStorager // Хранилище статистики сервера

	Config
	Deleter
}

// ServerBuilder абстракция для создания сервера
type ServerBuilder struct {
	server Server
}

// SetConfigFromFile заполняет конфигурацию данными конфигурационного файла
func (s ServerBuilder) SetConfigFromFile(configName string) ServerBuilder {
	if configName == "" {
		configName = "config.json"
	}
	file, err := os.Open(configName)
	if err != nil {
		log.Println("SetConfigFromFile", err)
		return s
	}
	defer file.Close()

	var (
		serverCfg ServerConfig
		tlsCfg    TLSConfig
		dbCfg     DatabaseConfig
	)

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Ошибка при чтении файла", err)
		return s
	}

	err = json.Unmarshal(data, &serverCfg)
	if err != nil {
		log.Println("SetConfigFromFile 1", err)
	}
	err = json.Unmarshal(data, &tlsCfg)
	if err != nil {
		log.Println("SetConfigFromFile 2", err)
	}
	err = json.Unmarshal(data, &dbCfg)
	if err != nil {
		log.Println("SetConfigFromFile 3", err)
	}

	cfg := Config{ServerConfig: serverCfg, DatabaseConfig: dbCfg, TLSConfig: tlsCfg}

	log.Println("Конфиг из файла", cfg)

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

	if s.server.Config.TrustedSubnet == "" {
		s.server.Config.TrustedSubnet = cfg.TrustedSubnet
	}

	s.server.Config.GRPCAddress = cfg.GRPCAddress

	log.Println(s.server.Config)
	return s
}

// SetConfig заполняет конфигурацию данными из переменных окружения и флагов
func (s ServerBuilder) SetConfig(sa, bu, fsp, dn, eh, tn string) ServerBuilder {
	s.server.Config = Config{
		ServerConfig: ServerConfig{
			ServerAddress:   sa,
			BaseURL:         bu,
			FileStoragePath: fsp,
			TrustedSubnet:   tn,
		},
		DatabaseConfig: DatabaseConfig{DatabaseDSN: dn},
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

func (s ServerBuilder) SetTimeOuts(readTimeOut, writeTimeOut, shotDownTimeOut time.Duration) ServerBuilder {
	s.server.Config.ShutdownTimeout = shotDownTimeOut
	s.server.Config.ServerReadTimeout = readTimeOut
	s.server.Config.ServerWriteTimeout = writeTimeOut
	return s
}

// Build возвращает сконфигурированный сервер объект Car.
func (s ServerBuilder) Build() Server {
	return s.server
}

// SetStorages выбирает реализацию каждого интерфейса из трех. Костыльная зависимость от конкертного
func (s *Server) SetStorages(url URLStorager, user UserStorager, db DBStorager, st StatStorager) {
	s.URLStorage = url
	s.UserStorage = user
	s.DBStorage = db
	s.StatStorage = st
}
