// Package api - это пакет для работы с серверной частью приложения
package api

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"

	"github.com/rodeorm/shortener/internal/api/middleware"
)

// Server веб-сервер и его характеристики
type Server struct {
	ProfileType int // Тип профилирования (если необходимо)

	URLStorage  URLStorager  // Хранилище данных для URL
	UserStorage UserStorager // Хранилище данных для URL
	DBStorage   DBStorager   // Хранилище данных для DB

	Config
	Deleter
}

// ServerStart запускает веб-сервер
func ServerStart(s *Server) error {

	if s.URLStorage == nil || s.UserStorage == nil {
		return fmt.Errorf("не определены хранилища")
	}

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

	if s.Config.EnableHTTPS {
		m := newTLSManager(s.Config.ServerAddress)
		srv.TLSConfig = m.TLSConfig()
		err := srv.ListenAndServeTLS("", "")
		if err != nil {
			return err
		}
	} else {
		err := srv.ListenAndServe()
		if err != nil {
			return err
		}
	}

	return nil
}
