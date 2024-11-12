package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

/*
APIUserDeleteURLsHandler - это хэндлер для метода DELETE /api/user/urls, который в теле запроса принимает список идентификаторов сокращённых URL для асинхронного удаления. Запрос может быть таким:

Хендлер DELETE /api/user/urls, который в теле запроса принимает список идентификаторов сокращённых URL для асинхронного удаления. Запрос может быть таким:
*/
func (h *Server) APIUserDeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	w, user, err := h.getUserIdentity(w, r)
	if user.WasUnathorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		handleError(w, err, "APIUserDeleteURLsHandler 1")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, "APIUserDeleteURLsHandler 2")
		return
	}

	// Помещаем URL в очередь на асинхронное удаление. В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted.
	urls, err := core.GetURLsFromString(string(bodyBytes), user)

	if err != nil {
		handleError(w, err, "APIUserDeleteURLsHandler 3")
		return
	}
	err = h.DeleteQueue.Push(urls)
	if err != nil {
		handleError(w, err, "APIUserDeleteURLsHandler 4")
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, string(bodyBytes))
}
