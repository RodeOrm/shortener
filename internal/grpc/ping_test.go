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

	// Инициализируем gRPC сервер с мокированным хранилищем
	grpcSrv := grpcServer{Server: core.Server{DBStorage: storage}}
	grpcSrv.srv = grpc.NewServer(grpc.UnaryInterceptor(interc.UnaryLogInterceptor))

	pb.RegisterURLServiceServer(grpcSrv.srv, &grpcSrv)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		lis, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcSrv.srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Создаём gRPC клиент для отправки запросов
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewURLServiceClient(conn)

	// Подготавливаем контекст и выполняем вызов PingDB
	ctx := context.Background()
	var header metadata.MD

	pingDBResponse, err := c.PingDB(ctx, &pb.PingDBRequest{}, grpc.Header(&header))
	if err != nil {
		log.Println("Ошибка при вызове PingDB:", err)
		t.FailNow() // Завершаем тест с ошибкой, если вызов не удался
	}
	st, _ := status.FromError(err)
	log.Printf("Результаты PingDB: %v", pingDBResponse)

	assert.Equal(t, codes.OK, st.Code())
}
