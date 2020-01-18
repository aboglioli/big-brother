package workers

import "time"

type Task interface {
	Done() <-chan bool
	Err() <-chan error
}

type task struct {
	fn   func() error
	done chan bool
	err  chan error
}

func (t *task) Done() <-chan bool {
	return t.done
}

func (t *task) Err() <-chan error {
	return t.err
}

// Interface
type Queue interface {
	Do(fn func() error) Task
	Run()
	Finish()
}

type queue struct {
	workers int
	tasks   chan *task
	retries int
	sleep   time.Duration
}

func NewQueue(workers int, retries int, sleepBetweenCalls time.Duration) Queue {
	return &queue{
		workers: workers,
		tasks:   make(chan *task, 16),
		retries: retries,
		sleep:   sleepBetweenCalls,
	}
}

func (q *queue) Do(fn func() error) Task {
	task := &task{
		fn:   fn,
		done: make(chan bool),
		err:  make(chan error),
	}

	go func() {
		q.tasks <- task
	}()

	return task
}

func (q *queue) Run() {
	for i := 0; i < q.workers; i++ {
		go q.worker()
	}
}

func (q *queue) Finish() {
	close(q.tasks)
}

func (q *queue) worker() {
	for t := range q.tasks {
		q.runTask(t)
	}
}

func (q *queue) runTask(t *task) {
	defer func() {
		close(t.done)
		close(t.err)
	}()

	for i := 0; i < q.retries; i++ {
		err := t.fn()
		if err == nil {
			select {
			case t.done <- true:
			default:
			}
			return
		}

		select {
		case t.err <- err:
			time.Sleep(q.sleep)
		case <-time.After(q.sleep):
		}
	}

	select {
	case t.done <- false:
	default:
	}
}
