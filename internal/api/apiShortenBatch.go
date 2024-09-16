package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

func (h Server) APIShortenBatch(w http.ResponseWriter, r *http.Request) {
	w, user, _, err := h.GetUserIdentity(w, r)

	if err != nil {
		log.Println("APIShortenBatch", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		urlReq []core.URLWithCorrelationRequest
		urlRes []core.URLWithCorrelationResponse
	)

	bodyBytes, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(bodyBytes, &urlReq)

	if err != nil {
		log.Println("APIShortenBatch", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, value := range urlReq {
		shortURLKey, _, err := h.Storage.InsertURL(value.Origin, h.BaseURL, user)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		urlResPart := core.URLWithCorrelationResponse{CorID: value.CorID, Short: h.BaseURL + "/" + shortURLKey}
		urlRes = append(urlRes, urlResPart)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	bodyBytes, err = json.Marshal(urlRes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
