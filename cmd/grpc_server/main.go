package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rodeorm/shortener/internal/grpc/server"
)

func main() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer()
	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
