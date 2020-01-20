package workers

import (
	"context"
	"time"
)

type Result struct {
	Done    bool
	Err     error
	Success bool
}

type Queue interface {
	Do(fn func() error) <-chan Result
	Run()
	Finish()
}

type task struct {
	fn     func() error
	result chan Result
	ctx    context.Context
	cancel context.CancelFunc
}

type queue struct {
	workers int
	retries int
	sleep   time.Duration

	chTasks chan *task
	tasks   []*task
}

func NewQueue(workers int, retries int, sleepBetweenCalls time.Duration) Queue {
	return &queue{
		workers: workers,
		retries: retries,
		sleep:   sleepBetweenCalls,

		chTasks: make(chan *task),
		tasks:   make([]*task, 0),
	}
}
func (q *queue) Run() {
	for i := 0; i < q.workers; i++ {
		go q.worker()
	}
}

func (q *queue) Do(fn func() error) <-chan Result {
	ctx, cancel := context.WithCancel(context.Background())
	t := &task{
		fn:     fn,
		result: make(chan Result),
		ctx:    ctx,
		cancel: cancel,
	}
	q.addTask(t)
	return t.result
}

func (q *queue) Finish() {
	close(q.chTasks)
}

func (q *queue) addTask(t *task) {
	if len(q.tasks) == q.workers {
		fTask, rTasks := q.tasks[0], q.tasks[1:]
		fTask.cancel()
		q.tasks = rTasks
	}

	q.tasks = append(q.tasks, t)
	q.chTasks <- t
}

func (q *queue) worker() {
	for t := range q.chTasks {
		q.runTask(t)
	}
}

func (q *queue) runTask(t *task) {
	defer close(t.result)

	res := Result{
		Done:    false,
		Err:     nil,
		Success: false,
	}
	for i := 1; i <= q.retries; i++ {
		err := t.fn()
		if err == nil {
			res.Done = true
			res.Err = nil
			res.Success = true
			break
		}

		res.Err = err

		if i == q.retries {
			res.Done = true
		}

		select {
		case t.result <- res:
			if res.Done {
				return
			}
			time.Sleep(q.sleep)
		case <-time.After(q.sleep):
		}
	}

	select {
	case t.result <- res:
	case <-t.ctx.Done():
	}
}
