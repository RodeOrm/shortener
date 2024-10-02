package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

/*APIUserGetURLsHandler возвращает пользователю все когда-либо сокращённые им URL в формате JSON*/
func (h Server) APIUserGetURLsHandler(w http.ResponseWriter, r *http.Request) {
	w, user, isUnathorized, err := h.GetUserIdentity(w, r)

	if isUnathorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		/*
			log.Println("APIUserGetURLsHandler 1", err)
			w.WriteHeader(http.StatusNoContent) */
		handleError(w, err, "APIUserGetURLsHandler 1")
		return
	}
	URLHistory, err := h.Storage.SelectUserURLHistory(user)
	if err != nil {
		handleError(w, err, "APIUserGetURLsHandler 1")
		return
	}

	//Не очень изящно, конечно. Т.к. не хочется слишком много мест переделывать
	history := make([]core.UserURLPair, 0)

	for _, v := range *URLHistory {
		pair := core.UserURLPair{UserKey: v.UserKey, Short: fmt.Sprintf("%s/%s", h.BaseURL, v.Short), Origin: v.Origin}
		history = append(history, pair)
	}

	bodyBytes, err := json.Marshal(history)
	if err != nil {
		handleError(w, err, "APIUserGetURLsHandler 1")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bodyBytes))
}
