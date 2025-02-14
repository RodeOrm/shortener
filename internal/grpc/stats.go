package grpc

import (
	"context"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/grpc/meta"
	pb "github.com/rodeorm/shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stats аналог хэндлера для api
func (g *grpcServer) Stats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	ip, err := meta.GetIPFromCtx(&ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, `IP не корректный`)
	}

	if !core.IsIPInCIDR(ip, g.TrustedSubnet) {
		return nil, status.Error(codes.PermissionDenied, `переданный IP не находится в CIDR`)
	}

	s, err := g.StatStorage.SelectStatistic()
	if err != nil {
		return nil, status.Error(codes.NotFound, `нет статистики или внутренняя ошибка`)
	}

	resp := pb.StatsResponse{Statistic: &pb.Statistic{Urls: int32(s.UrlQty), Users: int32(s.UsrQty)}}

	return &resp, nil
}
