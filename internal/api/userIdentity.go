package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
	"github.com/rodeorm/shortener/internal/core"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые, но возвращает ошибку
func (h Server) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *core.User, bool, error) {
	userKey, _ := cookie.GetUserKeyFromCoockie(r)
	var isUnathorized bool
	key, err := strconv.Atoi(userKey)
	if err != nil {
		isUnathorized = true
	}
	user, err := h.Storage.InsertUser(key)
	if err != nil {
		log.Fatal(err)
		return w, nil, isUnathorized, err
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(fmt.Sprint(user.Key)))
	return w, user, isUnathorized, err
}
