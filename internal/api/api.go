// Package api - это пакет для работы с серверной частью приложения
package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/rodeorm/shortener/internal/api/middleware"
	"github.com/rodeorm/shortener/internal/logger"
)

const (
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 15 * time.Second
	shutdownTimeout    = 30 * time.Second
)

// Server веб-сервер и его характеристики
type Server struct {
	srv             *http.Server  // Сервер
	idleConnsClosed chan struct{} // Уведомление о завершении работы

	ProfileType int // Тип профилирования (если необходимо)

	URLStorage    URLStorager    // Хранилище данных для URL
	UserStorage   UserStorager   // Хранилище данных для URL
	DBStorage     DBStorager     // Хранилище данных для DB
	ServerStorage ServerStorager // Хранилище статистики сервера

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
	r.HandleFunc("/api/internal/stats", s.APIStatsHandler).Methods(http.MethodGet)

	r.HandleFunc("/", s.badRequestHandler)
	r.Use(middleware.WithZip, middleware.WithLog)

	pprofRouter := r.PathPrefix("/debug/pprof/").Subrouter()
	pprofRouter.HandleFunc("/", pprof.Index)
	pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
	pprofRouter.HandleFunc("/profile", pprof.Profile)
	pprofRouter.HandleFunc("/symbol", pprof.Symbol)
	pprofRouter.HandleFunc("/trace", pprof.Trace)

	s.srv = &http.Server{
		Handler:      r,
		Addr:         s.ServerAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.gracefulShutDown()

	for i := 0; i < s.WorkerCount; i++ {
		w := NewWorker(i, s.DeleteQueue, s.URLStorage, s.BatchSize)
		go w.delete(s.idleConnsClosed)
	}

	if s.Config.EnableHTTPS {
		m := newTLSManager(s.Config.ServerAddress)
		s.srv.TLSConfig = m.TLSConfig()
		err := s.srv.ListenAndServeTLS("", "")
		// ждём завершения процедуры graceful shutdown
		<-s.idleConnsClosed
		// получили оповещение о завершении
		logger.Log.Info("Server Shutdowned",
			zap.String("Server Shutdowned gracefully", s.ServerAddress),
		)

		if err != nil {
			return err
		}
	} else {
		err := s.srv.ListenAndServe()
		// ждём завершения процедуры graceful shutdown
		<-s.idleConnsClosed
		// получили оповещение о завершении
		// например закрыть соединение с базой данных,
		// закрыть открытые файлы
		logger.Log.Info("Server Shutdowned",
			zap.String("Завершили изящное выключение", s.ServerAddress),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Server) gracefulShutDown() {
	// через этот канал сообщим основному потоку, что соединения закрыты
	h.idleConnsClosed = make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGQUIT)
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint

		// создаем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := h.srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			logger.Log.Error("Server Shutdowned",
				zap.String("Ошибка при изящном выключении", "Сервер без https"),
			)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		logger.Log.Info("Server Shutdown",
			zap.String("Начали изящное выключение", h.ServerAddress),
		)
		close(h.idleConnsClosed)
	}()
}
