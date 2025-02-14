// Package zip определяет реализацию сжатия/разжимания gzip
package zip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// GzipWriter - абстракция над Writer и ResponseWriter
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write - синтаксическое упрощение для доступа к методу io.Writer
func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// DecompressGzip осуществляет декомпрессию данных, сжатых gzip
func DecompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("ошибка при декомпрессии данных из gzip: %v", err)
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декомпрессии данных из gzip: %v", err)
	}

	return b.Bytes(), nil
}

// IsGzip  проверяет по заголовкам, поддерживается ли сжатие gzip
func IsGzip(headers map[string][]string) bool {
	for _, value := range headers["Content-Encoding"] {
		if value == "application/gzip" || value == "gzip" {
			return true
		}
	}
	return false
}
