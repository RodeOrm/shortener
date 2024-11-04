package repo

import (
	"testing"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	ms, _, _ := GetStorages("", "")

	tests := []struct {
		name    string
		URL     core.URL
		user    core.User
		wantErr bool
	}{

		{
			name: "Проверка InsertURL", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка SelectOriginalURL", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка InsertUser", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка insertUserURLPair", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка SelectUserByKey", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка SelectUserURLHistory", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка getNextFreeKey", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
		{
			name: "Проверка DeleteURLs", URL: core.URL{OriginalURL: "https://www.yandex.ru", Key: "1", UserKey: 2, HasBeenShorted: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Проверка InsertURL":
				_, err := ms.InsertURL(tt.URL.OriginalURL, tt.URL.OriginalURL, &tt.user)
				require.NoError(t, err)
			case "Проверка SelectOriginalURL":
				ms.SelectOriginalURL(tt.URL.OriginalURL)
			case "Проверка InsertUser":
				ms.InsertUser(tt.user.Key)
			case "Проверка insertUserURLPair":
				ms.insertUserURLPair(tt.URL.Key, tt.URL.OriginalURL, &tt.user)
			case "Проверка SelectUserByKey":
				ms.SelectUserByKey(tt.URL.UserKey)
			case "Проверка SelectUserURLHistory":
				ms.SelectUserURLHistory(&tt.user)
			case "Проверка getNextFreeKey":
				ms.getNextFreeKey()
			case "Проверка DeleteURLs":
				ms.SelectOriginalURL(tt.URL.Key)
			}

		})
	}
}
