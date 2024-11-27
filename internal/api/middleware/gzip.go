// Package middleware предназначен для работы с middleware
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/rodeorm/shortener/internal/zip"
)

// WithZip - middleware для сжатия/распаковки
func WithZip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		bodyBytes, _ := io.ReadAll(r.Body)
		if zip.IsGzip(r.Header) {
			bodyBytes, _ = zip.DecompressGzip(bodyBytes)
		}

		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(zip.GzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
