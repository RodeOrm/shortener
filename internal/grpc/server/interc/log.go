package interc

import (
	"context"
	"time"

	"github.com/rodeorm/shortener/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLogInterceptor перехватчик для логирования
func UnaryLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Записываем начало времени
	startTime := time.Now()

	// Вызываем следующий обработчик
	resp, err := handler(ctx, req)
	// Вычисляем продолжительность
	duration := time.Since(startTime)

	// Получаем статус
	var grpcStatus *status.Status
	if err != nil {
		grpcStatus = status.Convert(err)
	} else {
		grpcStatus = status.New(0, "OK")
	}

	// Логируем информацию
	logger.Log.Info("gRPC call in log-interceptor",
		zap.String("method", info.FullMethod),
		zap.Duration("duration", duration),
		zap.String("status", grpcStatus.Message()),
		zap.Int64("size", int64(len(req.(string)))), // Считаем размер запроса как длину строки запроса
	)

	return resp, err
}
