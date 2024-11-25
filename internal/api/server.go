package api

import (
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

// Server веб-сервер и его характеристики
type httpServer struct {
	srv *http.Server
	core.Server
}
