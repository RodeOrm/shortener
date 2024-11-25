package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/rodeorm/shortener/proto"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewURLServiceClient(conn)
	a, err := c.PingDB(context.Background(), &pb.PingDBRequest{})
	fmt.Println("Hello, from grpc client", a, err)
}
