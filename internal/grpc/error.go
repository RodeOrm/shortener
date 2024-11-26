package grpc

import (
	"time"

	"github.com/rodeorm/shortener/internal/logger"
	"go.uber.org/zap"
)

func handleError(f string, err error) {
	logger.Log.Info(f,
		zap.String(err.Error(), time.Now().String()),
	)
}
