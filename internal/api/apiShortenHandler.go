package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

// APIShortenHandler принимает в теле запроса JSON-объект {"url":"<some_url>"} и возвращает в ответ объект {"result":"<shorten_url>"}.
func (h Server) APIShortenHandler(w http.ResponseWriter, r *http.Request) {
	url := core.URL{}
	shortURL := core.ShortenURL{}

	w, user, _, err := h.GetUserIdentity(w, r)
	if err != nil {
		log.Println("APIShortenHandler 1", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(bodyBytes, &url)
	if err != nil {
		log.Println("APIShortenHandler 2", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURLKey, isDuplicated, err := h.Storage.InsertURL(url.Key, h.BaseURL, user)
	if err != nil {
		log.Println("APIShortenHandler 3", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	shortURL.Key = h.BaseURL + "/" + shortURLKey
	if isDuplicated {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	bodyBytes, err = json.Marshal(shortURL)
	if err != nil {
		log.Println("APIShortenHandler 4", err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
