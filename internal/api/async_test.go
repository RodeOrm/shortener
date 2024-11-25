package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/mocks"
	"github.com/stretchr/testify/require"
)

func TestAsync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mocks.NewMockStorager(ctrl)
	storage.EXPECT().DeleteURLs(gomock.Any()).Return(nil).AnyTimes()

	q := core.NewQueue(1)
	require.NotNil(t, q)

	w := core.NewWorker(1, q, storage, 1)
	require.NotNil(t, w)
	close := make(chan struct{})
	go w.Delete(close)

	url := []core.URL{
		{OriginalURL: "http://yandex.ru"},
	}
	err := q.Push(url)
	require.NoError(t, err)
}

func TestNewQueue(t *testing.T) {
	require.NotNil(t, core.NewQueue(10))
}

func TestPush(t *testing.T) {
	queue := core.NewQueue(10)
	require.NotNil(t, queue)
	urls := []core.URL{{OriginalURL: "https://yandex.ru"}, {OriginalURL: "https://yandex.com"}}
	err := queue.Push(urls)
	require.NoError(t, err)
}
