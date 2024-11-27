package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/grpc/interc"
	"github.com/rodeorm/shortener/mocks"
	pb "github.com/rodeorm/shortener/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mocks.NewMockStatStorager(ctrl)
	storage.EXPECT().SelectStatistic().Return(&core.ServerStatistic{UrlQty: 100, UsrQty: 10}, nil).AnyTimes()

	grpcSrv := grpcServer{Server: core.Server{StatStorage: storage, Config: core.Config{ServerConfig: core.ServerConfig{TrustedSubnet: "10.0.0.0/24"}}}}
	grpcSrv.srv = grpc.NewServer(grpc.UnaryInterceptor(interc.UnaryLogInterceptor))
	defer grpcSrv.srv.Stop()

	pb.RegisterURLServiceServer(grpcSrv.srv, &grpcSrv)

	go func() {
		lis, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcSrv.srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()

	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewURLServiceClient(conn)

	// добавляем мету с IP к запросу
	md := metadata.New(map[string]string{"x-real-ip": "10.0.0.0"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var header metadata.MD

	resp, err := c.Stats(ctx, &pb.StatsRequest{}, grpc.Header(&header))
	if err != nil {
		log.Println("Ошибка при вызове Stats:", err)
		t.FailNow()
	}
	st, _ := status.FromError(err)
	log.Printf("Результаты StatDB: %v", resp)

	assert.Equal(t, codes.OK, st.Code(), "сервер возвращает некорректный код")
	assert.Equal(t, resp.Statistic.Urls, int32(100), "сервер возвращает некорректное количество урл")
	assert.Equal(t, resp.Statistic.Users, int32(10), "сервер возвращает некорректное количество пользователей")
}
