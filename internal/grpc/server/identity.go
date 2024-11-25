package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/grpc/server/meta"
)

// getUserIdentity определяет по контексту какой пользователь авторизовался, если мета в контексте некорректная, то создает нового пользователя и новую мету в контекст,
// возвращает совместно с ними и ошибку
func (g *grpcServer) getUserIdentity(ctx *context.Context) (*context.Context, error) {
	userKey, err := meta.GetUserKeyFromCtx(ctx)
	user := &core.User{}

	if err != nil {
		user.WasUnathorized = true
	}

	key, err := strconv.Atoi(userKey)
	// Если идентификатор - это не число, то пользователь точно не авторизован. key остается со значением по умолчанию.
	if err != nil {
		user.WasUnathorized = true
	}
	user, err = g.UserStorage.InsertUser(key)
	if err != nil {
		return nil, err
	}
	ctx, err = meta.PutUserKeyToCtx(fmt.Sprint(user.Key))
	if err != nil {
		return nil, err
	}
	return ctx, nil
}
