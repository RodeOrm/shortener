package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rodeorm/shortener/internal/core"
)

/*APIUserGetURLsHandler возвращает пользователю все когда-либо сокращённые им URL в формате JSON*/
func (h Server) APIUserGetURLsHandler(w http.ResponseWriter, r *http.Request) {
	w, userKey := h.GetUserIdentity(w, r)
	userID, err := strconv.Atoi(userKey)
	if err != nil {
		fmt.Println("Проблемы с получением пользователя", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	URLHistory, err := h.Storage.SelectUserURLHistory(userID)

	history := make([]core.UserURLPair, 0)
	for _, v := range *URLHistory {
		pair := core.UserURLPair{UserKey: v.UserKey, Short: fmt.Sprintf("%s/%s", h.BaseURL, v.Short), Origin: v.Origin}
		history = append(history, pair)
	}

	if err != nil {
		fmt.Println("Проблемы с получением истории пользователя", err)
		w.WriteHeader(http.StatusNoContent)
		return
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
