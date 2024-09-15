package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

/*APIUserGetURLsHandler возвращает пользователю все когда-либо сокращённые им URL в формате JSON*/
func (h Server) APIUserGetURLsHandler(w http.ResponseWriter, r *http.Request) {
	w, user, err := h.GetUserIdentity(w, r)
	if err != nil {
		log.Println("APIUserGetURLsHandler", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	URLHistory, err := h.Storage.SelectUserURLHistory(user)
	if err != nil {
		fmt.Println("APIUserGetURLsHandler", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//Не очень изящно, конечно. Т.к. не хочется слишком много мест
	history := make([]core.UserURLPair, 0)

	for _, v := range *URLHistory {
		pair := core.UserURLPair{UserKey: v.UserKey, Short: fmt.Sprintf("%s/%s", h.BaseURL, v.Short), Origin: v.Origin}
		history = append(history, pair)
	}

	bodyBytes, err := json.Marshal(history)
	if err != nil {
		fmt.Println("Проблемы при маршалинге истории урл", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bodyBytes))
}
