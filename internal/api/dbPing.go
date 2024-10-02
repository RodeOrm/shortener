package api

import (
	"fmt"
	"net/http"
)

func (h Server) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Storage.PingDB()
	if err == nil {
		fmt.Fprintf(w, "%s", "Успешное соединение с БД")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Ошибка соединения с БД")
	}
}
