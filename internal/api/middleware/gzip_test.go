package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithZip(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	zipHandler := WithZip(testHandler)

	req := httptest.NewRequest(http.MethodPost, "http://example.com/test", io.NopCloser(strings.NewReader("Something")))
	req.Header.Add("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()
	zipHandler.ServeHTTP(rr, req)
	res := rr.Result()

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("Код ответа не соответствует ожидаемому. Тело запроса: %s", body))
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"), fmt.Sprintf("Header не содержит Content-Encoding gzip. Header: %v", res.Header))

}

func TestGzipWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	gzipW := gzip.NewWriter(&buf)
	writer := gzipWriter{Writer: gzipW}

	data := []byte("test data")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Не ожидаем ошибки, но получили %v", err)
	}

	err = gzipW.Close() // Закрываем gzip.Writer, чтобы записать данные в буфер
	if err != nil {
		t.Fatalf("не ожидаем ошибки при closing gzip writer, получили %v", err)
	}

	if n != len(data) {
		t.Errorf("ожидаем запись %d bytes, получили %d", len(data), n)
	}

	if buf.Len() == 0 {
		t.Error("ожидаем не нулевую длину сжатых данных")
	}
}
