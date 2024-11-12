package api

import (
	"fmt"
	"sync"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/logger"
)

// Worker структура, удаляющая URL
type Worker struct {
	id         int
	batchSize  int
	queue      *Queue
	urlStorage URLStorager
}

// NewWorker создает новый Worker
func NewWorker(id int, queue *Queue, storage URLStorager, batchSize int) *Worker {
	w := Worker{
		id:         id,
		queue:      queue,
		urlStorage: storage,
		batchSize:  batchSize,
	}
	return &w
}

// Push помещает пачку URL в очередь
func (q *Queue) Push(url []core.URL) error {
	var wg sync.WaitGroup

	for _, v := range url {
		wg.Add(1)
		go func() {
			q.ch <- &v
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

// NewQueue создает новую очередь URL размером n
func NewQueue(n int) *Queue {
	return &Queue{
		ch: make(chan *core.URL, n),
	}
}

// Queue очередь на удаление URL
type Queue struct {
	ch chan *core.URL
}

// PopWait извлекает пачку URL из очереди на удаление
func (q *Queue) popWait(n int) []core.URL {

	urls := make([]core.URL, 0)
	for i := 0; i < n; i++ {
		select {
		case val := <-q.ch:
			urls = append(urls, *val)
		default:
			continue
		}
	}
	return urls
}

// delete основной рабочий метод Worker, удаляющего url из очереди
func (w *Worker) delete(exit chan struct{}) {
	logger.Log.Info(fmt.Sprintf("воркер #%d стартовал", w.id))

	for {
		select {
		case _, ok := <-exit:
			if !ok {
				return
			}
		default:
			urls := w.queue.popWait(w.batchSize)

			if len(urls) == 0 {
				continue
			}
			err := w.urlStorage.DeleteURLs(urls)
			if err != nil {
				logger.Log.Error(fmt.Sprintf("ошибка при работе воркера %v стартовал", err))
				continue
			}
			logger.Log.Info(fmt.Sprintf("воркер #%d удалил пачку урл %v", w.id, urls))
		}
	}
}
