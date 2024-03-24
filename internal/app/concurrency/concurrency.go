package concurrency

import (
	"context"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app"
)

type Task struct {
	URL string
}

func NewQueue() *Queue {
	return &Queue{
		ch: make(chan *Task, 1),
	}
}

type Queue struct {
	ch chan *Task
}

func (q *Queue) Push(t *Task) {
	// добавляем задачу в очередь
	q.ch <- t
}

func (q *Queue) PopWait() *Task {
	// получаем задачу
	return <-q.ch
}

func NewWorker(workerId int, queue *Queue, deleter *Deleter) *Worker {
	w := Worker{
		workerId: workerId,
		queue:    queue,
		deleter:  deleter,
	}
	return &w
}

type Worker struct {
	workerId int
	queue    *Queue
	deleter  *Deleter
}

func (w *Worker) Loop() {
	for {
		t := w.queue.PopWait()

		err := w.deleter.Delete(t.URL)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		fmt.Printf("worker #%d deleted URL %s\n", w.workerId, t.URL)
	}
}

type Deleter struct {
	short   *app.URLShortenServiceMethods
	context context.Context
}

func NewDeleter(short *app.URLShortenServiceMethods, context context.Context) *Deleter {
	return &Deleter{
		short:   short,
		context: context,
	}
}

func (r *Deleter) Delete(url string) error {
	return (*r.short).DeleteLink(r.context, url)
}
