package meta

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func IsGzipSupported(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)

	return ok && md.Get("content-encoding") != nil && md.Get("content-encoding")[0] == "gzip"
}
