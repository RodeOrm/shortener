package api

import (
	"log"
	"sync"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/rodeorm/shortener/internal/repo"
)

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

func NewQueue(n int) *Queue {
	return &Queue{
		ch: make(chan *core.URL, n),
	}
}

type Queue struct {
	ch chan *core.URL
}

func (q *Queue) PopWait(n int) []core.URL {

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

type Worker struct {
	id        int
	batchSize int
	queue     *Queue
	storage   repo.AbstractStorage
}

func NewWorker(id int, queue *Queue, storage repo.AbstractStorage, batchSize int) *Worker {
	w := Worker{
		id:        id,
		queue:     queue,
		storage:   storage,
		batchSize: batchSize,
	}
	return &w
}

// Loop - основной рабочий метод воркера
func (w *Worker) Loop() {
	log.Println("Worker стартовал", w.id)
	for {
		urls := w.queue.PopWait(w.batchSize)

		if len(urls) == 0 {
			continue
		}
		err := w.storage.DeleteURLs(urls)
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}

		log.Printf("Воркер #%d удалил пачку урл %v", w.id, urls)
	}
}
