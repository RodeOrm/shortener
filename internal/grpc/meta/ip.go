package meta

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/metadata"
)

// GetIPFromCtx получает айпи из мета в контексте
func GetIPFromCtx(ctx *context.Context) (net.IP, error) {

	md, ok := metadata.FromIncomingContext(*ctx)
	if !ok {
		return nil, fmt.Errorf("не удалось извлечь метаданные")
	}

	// Получаем IP-адрес из метаданных (предположим, что клиент его нам предоставляет)
	ipStr := md.Get("x-real-ip")
	if len(ipStr) == 0 {
		ipStr = md.Get("x-forwarded-for")
	}

	if len(ipStr) == 0 {
		return nil, fmt.Errorf("IP-адрес не предоставлен")
	}

	ip := net.ParseIP(ipStr[0])
	if ip == nil {
		return nil, fmt.Errorf("IP-адрес не предоставлен")
	}

	return ip, nil
}
