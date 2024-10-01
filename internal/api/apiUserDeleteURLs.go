package api

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

/*
Добавьте в сервис новый хендлер DELETE /api/user/urls, который в теле запроса принимает список идентификаторов сокращённых URL для асинхронного удаления. Запрос может быть таким:
DELETE http://localhost:8080/api/user/urls
Content-Type: application/json

["6qxTVvsy", "RTfd56hn", "Jlfd67ds"]
В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted. Фактический результат удаления может происходить позже — оповещать пользователя об успешности или неуспешности не нужно.
Успешно удалить URL может пользователь, его создавший. При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone.
Совет:
Для эффективного проставления флага удаления в БД используйте множественное обновление (batch update).
Для максимального наполнения буфера объектов обновления используйте паттерн fanIn.
*/
func (h Server) APIUserDeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	w, user, isUnathorized, err := h.GetUserIdentity(w, r)
	if isUnathorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Println("APIUserDeleteURLsHandler 1", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("APIUserDeleteURLsHandler 2", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Помещаем URL в очередь на асинхронное удаление. В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted.
	urls, err := core.GetURLsFromString(string(bodyBytes), user)
	if err != nil {
		log.Println("APIUserDeleteURLsHandler 3", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.DeleteQueue.Push(urls)
	if err != nil {
		log.Println("APIUserDeleteURLsHandler 3", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// h.Storage.DeleteURLs(urls)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, string(bodyBytes))
}
