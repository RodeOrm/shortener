package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/rodeorm/shortener/internal/api/middleware"
	"github.com/rodeorm/shortener/internal/repo"
)

// ServerStart запускает веб-сервер
func ServerStart(s *Server) error {

	r := mux.NewRouter()
	r.HandleFunc("/", s.RootHandler).Methods(http.MethodPost)
	r.HandleFunc("/ping", s.PingDBHandler).Methods(http.MethodGet)
	r.HandleFunc("/{URL}", s.RootURLHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/shorten", s.APIShortenHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/user/urls", s.APIUserGetURLsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/user/urls", s.APIUserDeleteURLsHandler).Methods(http.MethodDelete)
	r.HandleFunc("/api/shorten/batch", s.APIShortenBatch).Methods(http.MethodPost)

	r.HandleFunc("/", s.BadRequestHandler)

	r.Use(middleware.ZipMiddleware, middleware.LogMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         s.ServerAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	s.Storage.CloseConnection()

	return nil
}

type Server struct {
	ServerAddress            string
	BaseURL                  string
	DatabaseConnectionString string

	Storage repo.AbstractStorage
}
