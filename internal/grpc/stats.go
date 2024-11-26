package grpc

import (
	"context"

	pb "github.com/rodeorm/shortener/proto"
)

// Stats аналог хэндлера для api
func (g *grpcServer) Stats(context.Context, *pb.StatsRequest) (*pb.StatsResponse, error) {
	var resp pb.StatsResponse

	return &resp, nil
}
