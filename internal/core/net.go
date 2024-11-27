package core

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// CheckNet проверяет, что IP-адрес клиента, переданный в заголовке запроса X-Real-IP, входит в доверенную подсеть
func CheckNet(r *http.Request, CIDR string) (bool, error) {
	/*
		Адрес доступен в поле структуры запроса Request.RemoteAddr из пакета net/http.
		Но есть проблема: если между пользователем и Go-приложением стоит HTTP-прокси,
		например nginx или envoy, то в этом поле получим адрес прокси, а не пользователя.
		Чтобы решить эту проблему, прокси-серверы нужно донастроить так, чтобы они прикладывали HTTP-заголовок к изначальному IP-адресу пользователя.
		Такой HTTP-заголовок обычно называют X-Real-IP, X-False-IP или X-Forwarded-For. Значение заголовка можно получить методом Request.Header.Get.
	*/

	// смотрим заголовок запроса X-Real-IP
	ipStr := r.Header.Get("X-Real-IP")
	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
		// этот заголовок содержит адреса отправителя и промежуточных прокси
		// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
		ips := r.Header.Get("X-Forwarded-For")
		// разделяем цепочку адресов
		ipStrs := strings.Split(ips, ",")
		// интересует только первый
		ipStr = ipStrs[0]
		// парсим
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return false, fmt.Errorf("ошибка при получении ip из http header")
	}

	inCIDR := IsIPInCIDR(ip, CIDR)
	return inCIDR, nil
}

// IsIPInCIDR проверяет, содержится ли IP в CIDR
func IsIPInCIDR(ip net.IP, cidrStr string) bool {
	_, cidr, err := net.ParseCIDR(cidrStr) // Парсинг CIDR из строки
	if err != nil {
		return false
	}

	return cidr.Contains(ip)
}
