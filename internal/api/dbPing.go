package api

import (
	"fmt"
	"net/http"

	"github.com/rodeorm/shortener/internal/logger"
)

// PingDBHandler - обработчик для метода GET для маршрута /ping, который при запросе проверяет соединение с базой данных.
//
// При успешной проверке возвращает статус 200 OK, при неуспешной — 500 Internal Server Error.
func (h *httpServer) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	err := h.DBStorage.Ping()
	if err == nil {
		logger.Log.Info("успешная попытка ping")
		fmt.Fprintf(w, "%s", "Успешное соединение с БД")
		w.WriteHeader(http.StatusOK)
	} else {
		logger.Log.Info("провальная попытка ping")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Ошибка соединения с БД")
	}
}
