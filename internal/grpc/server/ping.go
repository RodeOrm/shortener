package server

import (
	"context"
	"log"

	po "github.com/rodeorm/shortener/proto"
)

// PingDB пингует БД
func (Server) PingDB(context.Context, *po.PingDBRequest) (*po.PingDBResponse, error) {
	var resp po.PingDBResponse
	log.Println("Hello, from grpc server")
	return &resp, nil
}
