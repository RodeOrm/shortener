package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryLogInterceptor перехватчик. TODO!!!
func UnaryLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	/*
		var token string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("token")
			if len(values) > 0 {
				token = values[0]
			}
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}
	*/
	return handler(ctx, req)
}
