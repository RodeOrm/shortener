package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRootURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	storage.EXPECT().SelectOriginalURL(gomock.Any()).Return(&core.URL{Key: "Short", HasBeenDeleted: true}, nil).AnyTimes()

	cs := &core.Server{URLStorage: storage, UserStorage: storage}
	s := httpServer{Server: cs}

	handler := http.HandlerFunc(s.RootURLHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		requestURL   string
		expectedCode int
		expectedBody string
	}{
		{name: "проверка на попытку получить удаленный урл", method: http.MethodGet, expectedCode: http.StatusGone, requestURL: "deleted"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + "/" + tc.requestURL
			resp, err := req.Send()

			assert.NoError(t, err, fmt.Sprintf("ошибка при попытке сделать запрос %s, ошибка: %s", req.URL, err))
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Код ответа не соответствует ожидаемому")
		})
	}

}
