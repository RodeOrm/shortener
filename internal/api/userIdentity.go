package api

import (
	"fmt"
	"net/http"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые
func (h Server) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, string) {
	userKey, err := cookie.GetUserKeyFromCoockie(r)

	if err != nil {

		user, _ := h.Storage.InsertUser(0)
		userKey = fmt.Sprint(user.Key)
		http.SetCookie(w, cookie.PutUserKeyToCookie(userKey))
		return w, userKey
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(userKey))
	return w, userKey
}
