package api

import "golang.org/x/crypto/acme/autocert"

func newTLSManager(domain string) *autocert.Manager {
	return &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("certs"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist(domain),
	}
}
