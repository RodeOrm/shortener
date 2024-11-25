package api

import (
	"net/http"
)

// BadRequestHandler обрабатывает некорректные запросы и возвращает статус 400 BadRequest
func (h httpServer) badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

// ForbiddenHandler обрабатывает запросы, не имеющие доступа
func (h httpServer) ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
