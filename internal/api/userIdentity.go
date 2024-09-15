package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
	"github.com/rodeorm/shortener/internal/core"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые
func (h Server) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *core.User, error) {
	userKey, _ := cookie.GetUserKeyFromCoockie(r)
	key, _ := strconv.Atoi(userKey)
	user, err := h.Storage.InsertUser(key)
	if err != nil {
		log.Fatal(err)
		return w, nil, err
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(fmt.Sprint(user.Key)))
	return w, user, nil
}
