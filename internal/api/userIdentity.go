package api

import (
	"fmt"
	"net/http"
	"strconv"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
	"github.com/rodeorm/shortener/internal/core"
)

// getUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает нового пользователя и новые куки, но возвращает совместно с ними и ошибку
func (h Server) getUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *core.User, error) {
	userKey, err := cookie.GetUserKeyFromCookie(r)
	user := &core.User{}

	if err != nil {
		user.WasUnathorized = true
	}

	key, err := strconv.Atoi(userKey)
	// Если идентификатор - это не число, то пользователь точно не авторизован. key остается со значением по умолчанию.
	if err != nil {
		user.WasUnathorized = true
	}
	user, err = h.UserStorage.InsertUser(key)
	if err != nil {
		handleError(w, err, "GetUserIdentity")
	}
	cookie, err := cookie.PutUserKeyToCookie(fmt.Sprint(user.Key))
	if err != nil {
		return w, nil, err
	}
	http.SetCookie(w, cookie)
	return w, user, nil
}
