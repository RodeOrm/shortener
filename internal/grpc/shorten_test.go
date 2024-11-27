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
	grpcSrv := grpcServer{Server: &core.Server{URLStorage: repo.GetMemoryStorage(),
		UserStorage: repo.GetMemoryStorage(),
		Config: core.Config{
			ServerConfig: core.ServerConfig{BaseURL: "base.com"}}}}
	grpcSrv.srv = grpc.NewServer(grpc.UnaryInterceptor(interc.UnaryLogInterceptor))

	pb.RegisterURLServiceServer(grpcSrv.srv, &grpcSrv)
	defer grpcSrv.srv.Stop()

	go func() {
		lis, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcSrv.srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()

	type test struct {
		name     string
		want     codes.Code
		request  pb.ShortenRequest
		response pb.ShortenResponse
	}

	ts := test{
		name:     "Проверка обработки корректных запросов",
		want:     codes.OK,
		request:  pb.ShortenRequest{Url: "{\"url\":\"https://www.google.com\"}"},
		response: pb.ShortenResponse{},
	}
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewURLServiceClient(conn)

	ctx := context.Background()
	var header metadata.MD

	t.Run(ts.name, func(t *testing.T) {
		resp, err := c.Shorten(ctx, &ts.request, grpc.Header(&header))
		if err != nil {
			log.Println("Ошибка при вызове Shorten:", err)
			t.FailNow()
		}
		st, _ := status.FromError(err)
		log.Printf("Результаты Shorten: %v", resp.Url)

		assert.NoError(t, err, "ошибка при попытке сделать запрос")
		assert.Equal(t, ts.want, st.Code(), "Код ответа не соответствует ожидаемому")
	})

}
