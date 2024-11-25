// Package grpc используется для реализации gprc
package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/grpc/server/interc"
	pb "github.com/rodeorm/shortener/proto"
)

// Server поддерживает все необходимые методы сервера через встраивание pb.UnimplementedURLServiceServer
type grpcServer struct {
	pb.UnimplementedURLServiceServer
	core.Server
}

// ServerStart запускает grpc-сервер
func ServerStart(srv *core.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	listen, err := net.Listen("tcp", srv.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interc.UnaryLogInterceptor))
	pb.RegisterURLServiceServer(s, &grpcServer{Server: *srv})

	fmt.Println("Сервер gRPC начал работу на порту", srv.GRPCAddress)
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
