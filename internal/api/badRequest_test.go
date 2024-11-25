package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rodeorm/shortener/internal/core"
)

func TestBadRequestHandler(t *testing.T) {
	s := httpServer{Server: core.Server{
		Config: core.Config{ServerConfig: core.ServerConfig{BaseURL: "http:tiny.com"}}}}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.badRequestHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Хочу получить статус %v, получаю %v", http.StatusBadRequest, status)
	}
}
