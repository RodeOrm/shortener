// Package grpc используется для реализации gprc
package server

import (
	"context"

	"google.golang.org/grpc"

	"github.com/rodeorm/shortener/internal/grpc/server/interceptor"
	pb "github.com/rodeorm/shortener/proto"
)

// Server поддерживает все необходимые методы сервера через встраивание pb.UnimplementedURLServiceServer
type Server struct {
	pb.UnimplementedURLServiceServer
}

// NewServer создает новый инстанс сервера
func NewServer() *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryLogInterceptor))

	// регистрируем сервис
	pb.RegisterURLServiceServer(s, &Server{})

	return s
}

// Shorten аналог хэндлера Shorten для api
func (Server) Shorten(context.Context, *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var resp pb.ShortenResponse

	return &resp, nil
}

// GetUserURLs аналог хэндлера для api
func (Server) GetUserURLs(context.Context, *pb.UserURLsRequest) (*pb.UserURLsResponse, error) {
	var resp pb.UserURLsResponse

	return &resp, nil

}

// DeleteUserURLs аналог хэндлера для api
func (Server) DeleteUserURLs(context.Context, *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var resp pb.DeleteURLsResponse

	return &resp, nil

}

// Stats аналог хэндлера для api
func (Server) Stats(context.Context, *pb.StatsRequest) (*pb.StatsResponse, error) {
	var resp pb.StatsResponse

	return &resp, nil
}
