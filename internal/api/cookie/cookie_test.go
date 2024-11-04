package cookie

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	crypt "github.com/rodeorm/shortener/internal/crypt"
)

func TestGetUserKeyFromCookie(t *testing.T) {
	validKey := strconv.Itoa(100)
	invalidKey := "invalid"
	validToken, err := crypt.Encrypt(validKey)
	require.NoError(t, err)

	invalidToken, err := crypt.Encrypt(invalidKey)
	require.NoError(t, err)

	tests := []struct {
		cookie *http.Cookie

		name string
		Key  string

		err bool
	}{
		{name: "обработка валидных куки", Key: validKey, cookie: &http.Cookie{
			Name:   "token",
			Value:  validToken,
			MaxAge: 10000,
		}, err: false},
		{name: "обработка невалидных куки", Key: invalidKey, cookie: &http.Cookie{
			Name:   "token",
			Value:  invalidToken,
			MaxAge: 10000,
		}, err: true},
		{name: "обработка пустых куки", cookie: &http.Cookie{}, err: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := http.Request{Header: map[string][]string{
				"Accept-Encoding": {"gzip, deflate"},
				"Accept-Language": {"en-us"},
			}}
			req.AddCookie(tt.cookie)
			key, err := GetUserKeyFromCookie(&req)

			if tt.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.Key, key)
		})
	}
}

func TestPutUserKeyToCookie(t *testing.T) {
	validKey := strconv.Itoa(100)

	tests := []struct {
		name string
		Key  string
	}{
		{name: "обработка валидных идентификаторов", Key: validKey},
		{name: "обработка пустых идентификаторов"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie, err := PutUserKeyToCookie(tt.Key)
			require.NoError(t, err)
			require.NotNil(t, cookie)
		})
	}
}
