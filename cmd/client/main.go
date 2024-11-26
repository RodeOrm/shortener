package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/rodeorm/shortener/proto"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewURLServiceClient(conn)
	var header, trailer metadata.MD
	ctx := context.Background()
	pingDBResponse, err := c.PingDB(ctx, &pb.PingDBRequest{}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Println(err)
	}

	log.Println("Результаты PingDB", pingDBResponse, "Header: ", header)

	rootResponse, err := c.Root(ctx, &pb.RootRequest{Url: "http://www.yandex.ru"}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Println(err)
	}
	s, _ := status.FromError(err)

	log.Println("Результаты Root", rootResponse, "Header: ", header, s.Code())

	shortenRepsonse, err := c.Shorten(ctx, &pb.ShortenRequest{Url: "{\"url\":\"https://www.google.com\"}"}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Println(err)
	}
	s, _ = status.FromError(err)
	log.Println("Результаты Shorten", shortenRepsonse.Url, "Header: ", header, s.Code())
}
