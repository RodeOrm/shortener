package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAPIUserGetURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	userURLs := make([]core.UserURLPair, 10)
	userURLs = append(userURLs, core.UserURLPair{UserKey: 1000, Short: "1", Origin: "http://zzz.ru"})
	userURLs = append(userURLs, core.UserURLPair{UserKey: 1000, Short: "2", Origin: "http://yandex.com"})

	user := &core.User{Key: 1000}

	storage.EXPECT().InsertUser(gomock.Any()).Return(user, false, nil).AnyTimes()
	storage.EXPECT().SelectUserURLHistory(user).Return(userURLs, nil)

	s := Server{Storage: storage}

	handler := http.HandlerFunc(s.APIUserGetURLsHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{name: "проверка на попытку получить историю пользователя", method: http.MethodGet, expectedCode: http.StatusAccepted},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL

			resp, err := req.Send()

			assert.NoError(t, err, "ошибка при попытке сделать запрос", resp)
		})
	}

}
