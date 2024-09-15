package api

import (
	"fmt"
	"io"
	"net/http"
)

// RootHandler POST принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h Server) RootHandler(w http.ResponseWriter, r *http.Request) {

	w, userKey := h.GetUserIdentity(w, r)

	bodyBytes, _ := io.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	shortURLKey, isDuplicated, err := h.Storage.InsertURL(bodyString, h.BaseURL, userKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if isDuplicated {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	body := fmt.Sprintf("%s/%s", h.BaseURL, shortURLKey)
	w.Write([]byte(body))
}
