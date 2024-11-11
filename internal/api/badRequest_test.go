package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequestHandler(t *testing.T) {
	server := Server{}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.badRequestHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Хочу получить статус %v, получаю %v", http.StatusBadRequest, status)
	}
}
