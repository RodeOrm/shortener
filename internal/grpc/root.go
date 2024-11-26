package grpc

import (
	"context"

	"github.com/rodeorm/shortener/internal/core"
	pb "github.com/rodeorm/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Root аналог хэндлера Root
func (g *grpcServer) Root(ctx context.Context, req *pb.RootRequest) (*pb.RootResponse, error) {
	url := core.URL{OriginalURL: req.Url}

	user, md, err := g.getUserIdentity(&ctx)
	grpc.SetHeader(ctx, md)

	if err != nil {
		handleError("Grpc Root 1", err)
		return nil, err
	}
	u, err := g.Server.URLStorage.InsertURL(url.OriginalURL, g.BaseURL, user)
	if err != nil {
		handleError("Grpc Root 2", err)
		return nil, err
	}

	resp := pb.RootResponse{Shorten: g.BaseURL + "/" + u.Key}

	if url.HasBeenShorted {
		return &resp, status.Error(codes.AlreadyExists, `URL уже был сокращен`)
	}
	return &resp, status.Error(codes.OK, `URL принят к сокращению`)
}
