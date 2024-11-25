// Package meta
//
// Вместе с gRPC-запросами нужно передавать дополнительные данные — подобно передаче заголовков в HTTP-запросах.
// Метаданные хранятся в виде мапы, в которой ключу соответствует не одно значение, а слайс строк.
// Для передачи и получения метаданных используется context.Context. Тип MD служит для хранения метаданных.
package meta

import (
	"context"
	"strconv"

	"github.com/rodeorm/shortener/internal/crypt"
	"google.golang.org/grpc/metadata"
)

// GetUserKeyFromCtx получает идентфикатор пользователя из мета в контексте
func GetUserKeyFromCtx(ctx *context.Context) (string, error) {

	var token string

	md, ok := metadata.FromIncomingContext(*ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			token = values[0]
		}
	}

	userKey, err := crypt.Decrypt(token)
	if err != nil {
		return "", err
	}

	_, err = strconv.Atoi(userKey)

	if err != nil {
		return "", err
	}

	return userKey, nil
}

// PutUserKeyToCtx помещает идентификатор пользователя в контекст
func PutUserKeyToCtx(Key string) (*context.Context, error) {
	val, err := crypt.Encrypt(Key)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"token": val})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return &ctx, nil
}
