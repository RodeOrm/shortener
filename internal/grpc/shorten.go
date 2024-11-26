package grpc

import (
	"context"
	"encoding/json"

	"github.com/rodeorm/shortener/internal/core"
	pb "github.com/rodeorm/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Shorten аналог хэндлера Shorten для api
//
// принимает в теле запроса JSON-объект {"url":"<some_url>"} и возвращает в ответ объект {"result":"<shorten_url>"}.
func (g *grpcServer) Shorten(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var resp pb.ShortenResponse
	url := core.URL{}
	shortURL := core.ShortenURL{}

	user, md, err := g.getUserIdentity(&ctx)
	grpc.SetHeader(ctx, md)

	if err != nil {
		handleError("Grpc Shorten 1", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(req.Url), &url)
	if err != nil {
		handleError("Grpc Shorten 2", err)
		return nil, err
	}

	urlFromStorage, err := g.URLStorage.InsertURL(url.Key, g.BaseURL, user)
	url = *urlFromStorage

	if err != nil {
		handleError("Grpc Shorten 3", err)
		return nil, err
	}

	shortURL.Key = g.BaseURL + "/" + url.Key

	/*
		if url.HasBeenShorted {
			return &resp, status.Error(codes.AlreadyExists, `URL уже был сокращен`)
		}
	*/
	json, err := json.Marshal(shortURL)

	resp = pb.ShortenResponse{Url: string(json)}

	if err != nil {
		handleError("Grpc Shorten 3", err)
		return nil, err
	}

	return &resp, status.Error(codes.OK, `URL принят к сокращению`)
}
