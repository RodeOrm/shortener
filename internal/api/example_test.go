package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/repo"
	"github.com/rodeorm/shortener/mocks"
)

func ExampleServer_APIShortenHandler() {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name    string
		request string
		body    string

		server Server
		want   want
	}{
		{
			//Нужно принимать и возвращать JSON
			name:    "Проверка обработки корректных запросов: POST (json)",
			server:  Server{Config: Config{ServerAddress: "http://localhost:8080"}, URLStorage: repo.GetMemoryStorage(), UserStorage: repo.GetMemoryStorage()}, // С хранилищем в памяти, поэтому мокать  не надо
			body:    `{"url":"http://www.yandex.ru"}`,
			request: "http://localhost:8080/api/shorten",
			want:    want{statusCode: 201, contentType: "json"},
		},
		{
			//Нужно принимать и возвращать JSON
			name:    "Проверка обработки некорректных запросов: POST (json)",
			server:  Server{Config: Config{ServerAddress: "http://localhost:8080"}, URLStorage: repo.GetMemoryStorage(), UserStorage: repo.GetMemoryStorage()}, // С хранилищем в памяти, поэтому мокать  не надо
			body:    ``,
			request: "http://localhost:8080/api/shorten",
			want:    want{statusCode: 400},
		},
	}
	for _, tt := range tests {
		var request *http.Request
		if tt.body != "" {
			request = httptest.NewRequest(http.MethodPost, tt.request, bytes.NewReader([]byte(tt.body)))
		} else {
			request = httptest.NewRequest(http.MethodPost, tt.request, nil)
		}
		w := httptest.NewRecorder()
		h := http.HandlerFunc(tt.server.APIShortenHandler)
		h.ServeHTTP(w, request)
		result := w.Result()
		result.Body.Close()

	}
}

func ExampleServer_APIShortenBatchHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	storage.EXPECT().InsertUser(gomock.Any()).Return(&core.User{Key: 1000, WasUnathorized: false}, nil).MaxTimes(3)
	storage.EXPECT().InsertURL("http://double.com", gomock.Any(), gomock.Any()).Return(&core.URL{Key: "short", HasBeenShorted: true}, nil)
	storage.EXPECT().InsertURL("http://err", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("ошибка"))
	storage.EXPECT().InsertURL("http://valid.com", gomock.Any(), gomock.Any()).Return(&core.URL{Key: "short", HasBeenShorted: false}, nil)

	s := Server{UserStorage: storage, URLStorage: storage, DBStorage: storage}

	handler := http.HandlerFunc(s.APIShortenBatchHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	useCases := []struct {
		name         string
		method       string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{name: "попытка сократить ранее сокращенный урл", method: http.MethodPost, requestBody: "[" +
			"{" +
			"\"correlation_id\": \"<строковый идентификатор>\"," +
			"\"original_url\": \"http://double.com\"" +
			"}" +
			"]", expectedCode: http.StatusCreated}, // Непонятно, насколько правильно возвращать в этом случае "StatusCreated", т.к. другие хэндлеры для дублей возвращают конфликт
		{name: "проверка на попытку сократить невалидный урл", method: http.MethodPost, requestBody: "[" +
			"{" +
			"\"correlation_id\": \"<строковый идентификатор>\"," +
			"\"original_url\": \"http://err\"" +
			"}" +
			"]", expectedCode: http.StatusBadRequest},
		{name: "проверка на попытку сократить корректный урл, который не сокращали ранее", method: http.MethodPost, requestBody: "[" +
			"{" +
			"\"correlation_id\": \"<строковый идентификатор>\"," +
			"\"original_url\": \"http://valid.com\"" +
			"}" +
			"]", expectedCode: http.StatusCreated},
	}

	for _, uc := range useCases {
		req := resty.New().R()
		req.Method = uc.method
		req.URL = srv.URL
		req.Body = uc.requestBody

		resp, err := req.Send()
		if err == nil {
			fmt.Println(resp)
		}
	}
}

func ExampleServer_APIUserDeleteURLsHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	storage.EXPECT().InsertUser(gomock.Any()).Return(&core.User{Key: 1000, WasUnathorized: false}, nil).AnyTimes()
	storage.EXPECT().DeleteURLs(gomock.Any()).Return(nil).AnyTimes()

	s := Server{UserStorage: storage, URLStorage: storage, DBStorage: storage, Deleter: Deleter{DeleteQueue: NewQueue(3)}}

	handler := http.HandlerFunc(s.APIUserDeleteURLsHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	worker := NewWorker(1, s.DeleteQueue, storage, 1)
	go worker.delete()

	useCases := []struct {
		name         string
		method       string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{name: "проверка на попытку удалить банч урл", method: http.MethodPost, requestBody: "[" +
			"\"6qxTVvsy\", \"RTfd56hn\", \"Jlfd67ds\"" +
			"]", expectedCode: http.StatusAccepted},
	}

	for _, tc := range useCases {

		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL
		req.Body = tc.requestBody

		resp, err := req.Send()
		if err == nil {
			fmt.Println(resp)
		}
	}
}

func ExampleServer_APIUserGetURLsHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	userURLs := make([]core.UserURLPair, 0)
	userURLs = append(userURLs, core.UserURLPair{UserKey: 1000, Short: "1", Origin: "http://1.ru"})
	userURLs = append(userURLs, core.UserURLPair{UserKey: 1000, Short: "2", Origin: "http://2.com"})

	user := &core.User{Key: 1000, WasUnathorized: false}

	storage.EXPECT().InsertUser(gomock.Any()).Return(user, nil).AnyTimes()
	storage.EXPECT().SelectUserURLHistory(user).Return(userURLs, nil)

	s := Server{UserStorage: storage, URLStorage: storage, DBStorage: storage, Config: Config{BaseURL: "http:tiny.com"}}

	handler := http.HandlerFunc(s.APIUserGetURLsHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	useCases := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{name: "проверка на попытку получить историю пользователя", method: http.MethodGet, expectedCode: http.StatusAccepted},
	}

	for _, tc := range useCases {
		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL

		resp, err := req.Send()
		if err == nil {
			fmt.Println(resp)
		}
	}
}

func ExampleServer_PingDBHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	storage.EXPECT().Ping().Return(nil).AnyTimes()

	s := Server{UserStorage: storage, URLStorage: storage, DBStorage: storage}

	handler := http.HandlerFunc(s.PingDBHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	useCases := []struct {
		name         string
		method       string
		requestURL   string
		expectedCode int
		expectedBody string
	}{
		{name: "обработка успешной попытки достучаться к БД", method: http.MethodGet, expectedCode: http.StatusOK, requestURL: "/ping"},
	}

	for _, tc := range useCases {
		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL
		resp, err := req.Send()
		if err == nil {
			fmt.Println(resp)
		}
	}
}

func ExampleServer_RootHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	storage.EXPECT().InsertUser(gomock.Any()).Return(&core.User{Key: 1000, WasUnathorized: false}, nil).MaxTimes(3)
	storage.EXPECT().InsertURL("http://double.com", gomock.Any(), gomock.Any()).Return(&core.URL{Key: "short", HasBeenShorted: true}, nil)
	storage.EXPECT().InsertURL("http://err", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("ошибка"))
	storage.EXPECT().InsertURL("http://valid.com", gomock.Any(), gomock.Any()).Return(&core.URL{Key: "short", HasBeenShorted: false}, nil)

	s := Server{UserStorage: storage, URLStorage: storage, DBStorage: storage}

	handler := http.HandlerFunc(s.RootHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	useCases := []struct {
		name         string
		method       string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{name: "проверка на попытку сократить ранее сокращенный урл", method: http.MethodPost, requestBody: "http://double.com", expectedCode: http.StatusConflict},
		{name: "проверка на попытку сократить невалидный урл", method: http.MethodPost, requestBody: "http://err", expectedCode: http.StatusBadRequest},
		{name: "проверка на попытку сократить корректный урл, который не сокращали ранее", method: http.MethodPost, requestBody: "http://valid.com", expectedCode: http.StatusCreated},
	}

	for _, tc := range useCases {
		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL
		req.Body = tc.requestBody

		resp, err := req.Send()
		if err == nil {
			fmt.Println(resp)
		}
	}
}

func ExampleServer_RootURLHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)

	storage.EXPECT().SelectOriginalURL(gomock.Any()).Return(&core.URL{Key: "Short", HasBeenDeleted: false}, nil).AnyTimes()

	s := Server{UserStorage: storage, URLStorage: storage, DBStorage: storage}

	handler := http.HandlerFunc(s.RootURLHandler)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	useCases := []struct {
		name         string
		method       string
		requestURL   string
		expectedCode int
		expectedBody string
	}{
		{name: "редирект, если URL был сокращен ранее", method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect, requestURL: "http://www.yandex.ru"},
	}

	for _, tc := range useCases {
		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL + "/" + tc.requestURL
		resp, err := req.Send()
		if err == nil {
			fmt.Println(resp)
		}
	}
}
