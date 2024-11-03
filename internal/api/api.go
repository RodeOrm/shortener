// Package api - это пакет для работы с серверной частью приложения
package api

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"

	"github.com/rodeorm/shortener/internal/api/middleware"
)

// Server абстракция, отражающая веб-сервер и его характеристики
type Server struct {
	//Количество воркеров, асинхронно удаляющих url
	WorkerCount int
	//Размер пачки для удаления
	BatchSize int
	//Тип профилирования (если необходимо)
	ProfileType int

	//Адрес запуска веб-сервера
	ServerAddress string
	//Базовый URL для сокращенных адресов
	BaseURL string
	//Connection string для БД
	DatabaseConnectionString string

	//Хранилище данных для URL
	URLStorage URLStorager
	// Хранилище данных для URL
	UserStorage UserStorager
	// Хранилище данных для DB
	DBStorage DBStorager

	//Очередь удаления
	DeleteQueue *Queue
}

// ServerStart запускает веб-сервер
func ServerStart(s *Server) error {

	if s.DBStorage != nil {
		defer s.DBStorage.Close()
		defer close(s.DeleteQueue.ch)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", s.RootHandler).Methods(http.MethodPost)
	r.HandleFunc("/ping", s.PingDBHandler).Methods(http.MethodGet)
	r.HandleFunc("/{URL}", s.RootURLHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/shorten", s.APIShortenHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/user/urls", s.APIUserGetURLsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/user/urls", s.APIUserDeleteURLsHandler).Methods(http.MethodDelete)
	r.HandleFunc("/api/shorten/batch", s.APIShortenBatchHandler).Methods(http.MethodPost)

	r.HandleFunc("/", s.badRequestHandler)
	r.Use(middleware.WithZip, middleware.WithLog)

	pprofRouter := r.PathPrefix("/debug/pprof/").Subrouter()
	pprofRouter.HandleFunc("/", pprof.Index)
	pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
	pprofRouter.HandleFunc("/profile", pprof.Profile)
	pprofRouter.HandleFunc("/symbol", pprof.Symbol)
	pprofRouter.HandleFunc("/trace", pprof.Trace)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.ServerAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	for i := 0; i < s.WorkerCount; i++ {
		w := NewWorker(i, s.DeleteQueue, s.URLStorage, s.BatchSize)
		go w.delete()
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
