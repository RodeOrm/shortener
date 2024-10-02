package api

import (
	"fmt"
	"net/http"
	"strconv"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
	"github.com/rodeorm/shortener/internal/core"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые, но возвращает ошибку
func (h Server) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *core.User, bool, error) {
	userKey, err := cookie.GetUserKeyFromCoockie(r)
	var isUnathorized bool

	if err != nil {
		isUnathorized = true
	}

	key, err := strconv.Atoi(userKey)
	// Если идентификатор - это не число, то пользователь точно не авторизован
	if err != nil {
		isUnathorized = true
	}
	user, isUnathorized, err := h.Storage.InsertUser(key)
	if err != nil {
		handleError(w, err, "GetUserIdentity")
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(fmt.Sprint(user.Key)))
	return w, user, isUnathorized, nil
}
