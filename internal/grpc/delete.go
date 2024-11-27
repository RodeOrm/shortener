package grpc

import (
	"context"

	"github.com/rodeorm/shortener/internal/core"
	pb "github.com/rodeorm/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUserURLs помещает урлы в очередь на удаление
func (g *grpcServer) DeleteUserURLs(ctx context.Context, req *pb.DeleteURLsRequest) (*pb.DeleteURLsResponse, error) {
	var resp pb.DeleteURLsResponse
	user, md, err := g.getUserIdentity(&ctx)
	grpc.SetHeader(ctx, md)

	if user.WasUnathorized {
		return nil, status.Error(codes.Unauthenticated, `пользователь не был аутентифицирован`)
	}
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	urls, err := core.GetURLsFromString(req.UrlsToDelete, user)

	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	err = g.DeleteQueue.Push(urls)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	return &resp, status.Error(codes.OK, `URL добавлены в очередь на удаление`)
}
