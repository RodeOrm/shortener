package cookie

import (
	"fmt"
	"net/http"

	crypt "github.com/rodeorm/shortener/internal/crypt"
)

// GetUserKeyFromCookie получает идентфикатор пользователя из куки "токен"
func GetUserKeyFromCookie(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	if tokenCookie.Value == "" {
		return "", fmt.Errorf("не найдено актуальных cookie")
	}
	userKey, err := crypt.Decrypt(tokenCookie.Value)
	if err != nil {
		return "", err
	}
	return userKey, nil
}

// PutUserKeyToCookie помещает идентификатор пользователя в  куки "токен"
func PutUserKeyToCookie(Key string) *http.Cookie {
	val, _ := crypt.Encrypt(Key)

	cookie := &http.Cookie{
		Name:   "token",
		Value:  val,
		MaxAge: 10000,
	}
	return cookie
}
