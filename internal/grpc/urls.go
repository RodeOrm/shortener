package grpc

import (
	"context"

	pb "github.com/rodeorm/shortener/proto"
)

// GetUserURLs аналог хэндлера для api
func (g *grpcServer) GetUserURLs(context.Context, *pb.UserURLsRequest) (*pb.UserURLsResponse, error) {
	var resp pb.UserURLsResponse

	return &resp, nil

}
