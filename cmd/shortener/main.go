package main

import (
	"flag"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/rodeorm/shortener/internal/api"
	"github.com/rodeorm/shortener/internal/core"
	grpc "github.com/rodeorm/shortener/internal/grpc/server"
	"github.com/rodeorm/shortener/internal/repo"
)

/*
Сервис для сокращения длинных URL. Требования:
Сервер должен быть доступен по адресу: http://localhost:8080.
Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

func main() {
	flag.Parse()
	server, err := core.Configurate(a, b, c, config, d, f, w, s, q, p, bs, t)
	if err != nil {
		log.Fatal(err)
	}

	ms, fs, ps := repo.GetStorages(server.FileStoragePath, server.DatabaseDSN)
	if ps != nil {
		server.SetStorages(ps, ps, ps, ps)
	} else if fs != nil {
		server.SetStorages(fs, fs, nil, nil)
	} else if ms != nil {
		server.SetStorages(ms, ms, nil, nil)
	} else {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	core.Profile(server.ProfileType)

	go api.ServerStart(server, &wg)
	go grpc.ServerStart(server, &wg)

	wg.Wait()
}
