package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/grpc/interc"
	"github.com/rodeorm/shortener/internal/repo"
	pb "github.com/rodeorm/shortener/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestShortenServers(t *testing.T) {
	// Инициализируем gRPC сервер с мокированным хранилищем
	grpcSrv := grpcServer{Server: core.Server{URLStorage: repo.GetMemoryStorage(),
		UserStorage: repo.GetMemoryStorage(),
		Config: core.Config{
			ServerConfig: core.ServerConfig{BaseURL: "base.com"}}}}
	grpcSrv.srv = grpc.NewServer(grpc.UnaryInterceptor(interc.UnaryLogInterceptor))

	pb.RegisterURLServiceServer(grpcSrv.srv, &grpcSrv)
	defer grpcSrv.srv.Stop()

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

	tests := []struct {
		name     string
		want     codes.Code
		request  pb.ShortenRequest
		response pb.ShortenResponse
	}{

		{
			name:     "Проверка обработки корректных запросов",
			want:     codes.OK,
			request:  pb.ShortenRequest{Url: "{\"url\":\"https://www.google.com\"}"},
			response: pb.ShortenResponse{},
		},
	}
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewURLServiceClient(conn)

	// Подготавливаем контекст и выполняем вызов PingDB
	ctx := context.Background()
	var header metadata.MD

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Подготавливаем контекст и выполняем вызов PingDB
			resp, err := c.Shorten(ctx, &tc.request, grpc.Header(&header))
			if err != nil {
				log.Println("Ошибка при вызове Shorten:", err)
				t.FailNow() // Завершаем тест с ошибкой, если вызов не удался
			}
			st, _ := status.FromError(err)
			log.Printf("Результаты Shorten: %v", resp.Url)

			assert.NoError(t, err, "ошибка при попытке сделать запрос")
			assert.Equal(t, tc.want, st.Code(), "Код ответа не соответствует ожидаемому")
		})
	}
}
