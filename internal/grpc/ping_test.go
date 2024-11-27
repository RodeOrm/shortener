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
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/rodeorm/shortener/proto"
)

func TestDBPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)
	storage.EXPECT().Ping().Return(nil).AnyTimes()

	grpcSrv := grpcServer{Server: &core.Server{DBStorage: storage}}
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

	ctx := context.Background()
	var header metadata.MD

	pingDBResponse, err := c.PingDB(ctx, &pb.PingDBRequest{}, grpc.Header(&header))
	if err != nil {
		log.Println("Ошибка при вызове PingDB:", err)
		t.FailNow()
	}
	st, _ := status.FromError(err)
	log.Printf("Результаты PingDB: %v", pingDBResponse)

	assert.Equal(t, codes.OK, st.Code())
}
