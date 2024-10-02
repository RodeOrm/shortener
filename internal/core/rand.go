package core

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

// ReturnShortKey возвращает рандомный ключ (используем для генерации коротких URL)
func ReturnShortKey(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("некорректное значение ключа %v", n)
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}
