package server

import (
	"context"

	pb "github.com/rodeorm/shortener/proto"
)

// GetUserURLs аналог хэндлера для api
func (g *grpcServer) GetUserURLs(context.Context, *pb.UserURLsRequest) (*pb.UserURLsResponse, error) {
	var resp pb.UserURLsResponse

	return &resp, nil

}

// DeleteUserURLs аналог хэндлера для api
func (g *grpcServer) DeleteUserURLs(context.Context, *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var resp pb.DeleteURLsResponse

	return &resp, nil

}
