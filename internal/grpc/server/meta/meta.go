// Package meta
//
// Вместе с gRPC-запросами нужно передавать дополнительные данные — подобно передаче заголовков в HTTP-запросах.
// Метаданные хранятся в виде мапы, в которой ключу соответствует не одно значение, а слайс строк.
// Для передачи и получения метаданных используется context.Context. Тип MD служит для хранения метаданных.
package meta

import "context"

// GetUserKeyFromCtx получает идентфикатор пользователя из мета в контексте
func GetUserKeyFromCtx(ctx context.Context) (string, error) {
	return "", nil
}

// PutUserKeyToCtx помещает идентификатор пользователя в контекст
func PutUserKeyToCtx(Key string) *context.Context {
	return nil
}
