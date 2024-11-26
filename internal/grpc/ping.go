package grpc

import (
	"context"

	pb "github.com/rodeorm/shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PingDB пингует БД
func (g *grpcServer) PingDB(ctx context.Context, po *pb.PingDBRequest) (*pb.PingDBResponse, error) {
	var resp pb.PingDBResponse
	err := g.DBStorage.Ping()
	if err != nil {
		return nil, status.Error(codes.Unavailable, `ошибка при пинге БД`)
	}
	return &resp, status.Error(codes.OK, `успешный пинг`)
}
