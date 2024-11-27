// Package interc
package interc

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"github.com/rodeorm/shortener/internal/grpc/meta"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func GzipRequestInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if meta.IsGzipSupported(ctx) {

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)

		// Сериализуем запрос в байтовый массив
		data, err := proto.Marshal(req.(proto.Message))
		if err != nil {
			return err
		}

		// Сжимаем сериализованные данные
		if _, err := zw.Write(data); err != nil {
			return err
		}
		zw.Close()

		// Отправляем сжатый запрос
		return invoker(ctx, method, buf.Bytes(), reply, cc, opts...)
	}

	return invoker(ctx, method, req, reply, cc, opts...) // Отправляем без сжатия, если не указано
}

func GzipResponseInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, opts ...grpc.CallOption) error {
	if meta.IsGzipSupported(ctx) {
		// Разжимаем ответ
		gzReader, err := gzip.NewReader(bytes.NewReader(reply.([]byte)))
		if err != nil {
			return err
		}
		decompressedData, _ := io.ReadAll(gzReader)
		gzReader.Close()

		// Десериализуем данные обратно в структуру
		if err := proto.Unmarshal(decompressedData, reply.(proto.Message)); err != nil {
			return err
		}
	}

	return nil
}
