package server

import (
	"context"

	pb "github.com/rodeorm/shortener/proto"
)

// Shorten аналог хэндлера Shorten для api
func (g *grpcServer) Shorten(context.Context, *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var resp pb.ShortenResponse

	return &resp, nil
}
