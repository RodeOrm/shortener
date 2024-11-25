package core

import (
	"fmt"
	"sync"

	"github.com/rodeorm/shortener/internal/logger"
)

// Worker структура, удаляющая URL
type Worker struct {
	id         int
	batchSize  int
	queue      *Queue
	urlStorage URLStorager
}

// Queue очередь на удаление URL
type Queue struct {
	ch chan *URL
}

// StartWorkerPool один раз запускает воркер пул
func StartWorkerPool(workerCount int, deleteQueue *Queue, urlStorage URLStorager, batchSize int, idleClosedChan chan struct{}) {
	var once sync.Once
	once.Do(func() {
		for i := range workerCount {
			w := NewWorker(i, deleteQueue, urlStorage, batchSize)
			go w.Delete(idleClosedChan)
		}
	})
}

// Close закрывает очередь
func (q *Queue) Close() {
	close(q.ch)
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
func (q *Queue) Push(url []URL) error {
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
		ch: make(chan *URL, n),
	}
}

// PopWait извлекает пачку URL из очереди на удаление
func (q *Queue) PopWait(n int) []URL {

	urls := make([]URL, 0)
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

// Delete основной рабочий метод Worker, удаляющего url из очереди
func (w *Worker) Delete(exit chan struct{}) {
	logger.Log.Info(fmt.Sprintf("воркер #%d стартовал", w.id))

	for {
		select {
		case _, ok := <-exit:
			if !ok {
				return
			}
		default:
			urls := w.queue.PopWait(w.batchSize)

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
