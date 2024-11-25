// Package api - это пакет для работы с серверной частью приложения
package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/rodeorm/shortener/internal/api/middleware"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/logger"
)

// ServerStart запускает веб-сервер
func ServerStart(cs *core.Server, wg *sync.WaitGroup) error {
	defer wg.Done()

	if cs.URLStorage == nil || cs.UserStorage == nil {
		return fmt.Errorf("не определены хранилища")
	}

	if cs.DBStorage != nil {
		defer cs.DBStorage.Close()
		defer cs.DeleteQueue.Close()
	}

	s := httpServer{Server: *cs}

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
		WriteTimeout: s.ServerWriteTimeout,
		ReadTimeout:  s.ServerReadTimeout,
	}

	s.gracefulShutDown(wg)

	core.StartWorkerPool(s.WorkerCount, s.DeleteQueue, s.URLStorage, s.BatchSize, s.IdleConnsClosed)

	if s.Config.EnableHTTPS {
		m := newTLSManager(s.Config.ServerAddress)
		s.srv.TLSConfig = m.TLSConfig()
		err := s.srv.ListenAndServeTLS("", "")
		// ждём завершения процедуры graceful shutdown
		<-s.IdleConnsClosed
		// получили оповещение о завершении
		logger.Log.Info("https server shutdowned",
			zap.String("Завершили изящное выключение", s.ServerAddress),
		)

		if err != nil {
			return err
		}
	} else {
		err := s.srv.ListenAndServe()
		// ждём завершения процедуры graceful shutdown
		<-s.IdleConnsClosed
		// получили оповещение о завершении
		// например закрыть соединение с базой данных,
		// закрыть открытые файлы

		logger.Log.Info("http server shutdowned",
			zap.String("Завершили изящное выключение", s.ServerAddress),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *httpServer) gracefulShutDown(wg *sync.WaitGroup) {

	// через этот канал сообщим основному потоку, что соединения закрыты
	h.IdleConnsClosed = make(chan struct{})
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
		ctx, cancel := context.WithTimeout(context.Background(), h.ShutdownTimeout)
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
		close(h.IdleConnsClosed)
	}()
}
