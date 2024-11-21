package api

import "github.com/rodeorm/shortener/internal/core"

// DBStorager абстрация для методов, специфичных для БД
type DBStorager interface {

	// Close закрывает соединение
	Close()
	// Ping проверяет соединение
	Ping() error
}

// URLStorager абстрация для методов хранилища над URL
type URLStorager interface {

	//	InsertURL принимает оригинальный URL, базовый урл для генерации коротких адресов и пользователя.
	//
	//  Генерирует уникальный ключ для короткого адреса, сохраняет соответствие оригинального URL и ключа.
	//  Возвращает обновленный URL с соответствующим сокращенным URL, а также признаком того, что URL сократили ранее.
	InsertURL(URL, baseURL string, user *core.User) (*core.URL, error)

	// SelectOriginalURL возвращает URL на основании короткого
	SelectOriginalURL(shortURL string) (*core.URL, error)

	// DeleteURLs массово помечает URL как удаленные. Успешно удалить URL может только пользователь, его создавший.
	DeleteURLs(URLs []core.URL) error
}

// UserStorager абстрация для методов хранилища над User
type UserStorager interface {

	// InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии, а также значение "отсутствие авторизации по переданному идентификатору"
	InsertUser(Key int) (*core.User, error)

	// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
	SelectUserURLHistory(user *core.User) ([]core.UserURLPair, error)
}

type ServerStorager interface {
	SelectStatistic() (*core.ServerStatistic, error)
}
