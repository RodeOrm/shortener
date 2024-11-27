// Package grpc используется для реализации серверной части gprc
package grpc

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/grpc/interc"
	"github.com/rodeorm/shortener/internal/logger"
	pb "github.com/rodeorm/shortener/proto"
)

// GRPCServer поддерживает все необходимые методы сервера через встраивание pb.UnimplementedURLServiceServer,
// конфигурируется через атрибуты core.Server,
// включает в себя инстанс grpc.Server
type grpcServer struct {
	srv *grpc.Server
	*core.Server

	pb.UnimplementedURLServiceServer
}

// ServerStart запускает grpc-сервер
func ServerStart(core *core.Server, wg *sync.WaitGroup) {
	// Начинаем слушать порт из конфига
	listen, err := net.Listen("tcp", core.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpcServer{Server: core}
	grpcSrv.srv = grpc.NewServer(grpc.UnaryInterceptor(interc.UnaryLogInterceptor))

	pb.RegisterURLServiceServer(grpcSrv.srv, &grpcSrv)

	logger.Log.Info("grpc server started",
		zap.String("cервер gRPC начал работу на порту", grpcSrv.GRPCAddress),
	)

	go func() {
		if err := grpcSrv.srv.Serve(listen); err != nil {
			log.Fatalf("Ошибка при обработке: %v", err)
		}
	}()

	// Обработка сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	logger.Log.Info("grpc server gracefuly shutdown",
		zap.String("cервер gRPC начал изящно завершать работы на порту", grpcSrv.GRPCAddress),
	)
	grpcSrv.srv.GracefulStop() // Корректное завершение работы сервера
	logger.Log.Info("grpc server gracefuly shutdown",
		zap.String("cервер gRPC изящно завершил работу на порту", grpcSrv.GRPCAddress),
	)
	defer wg.Done()
}
