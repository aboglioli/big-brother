package workers

import "time"

type Task struct {
	Fn  func() error
	Err chan error
}

// Interface
type Queue interface {
	Do(fn func() error) <-chan error
	Run()
	Finish()
}

type queue struct {
	tasks   chan *Task
	retries int
	sleep   time.Duration
}

func NewQueue() Queue {
	return &queue{
		tasks:   make(chan *Task, 16),
		retries: 3,
		sleep:   500 * time.Millisecond,
	}
}

func (q *queue) Do(fn func() error) <-chan error {
	chErr := make(chan error)
	go func() {
		task := &Task{
			Fn:  fn,
			Err: chErr,
		}
		q.tasks <- task
	}()
	return chErr
}

func (q *queue) Run() {
	for task := range q.tasks {
		go q.runTask(task)
	}
}

func (q *queue) Finish() {
	close(q.tasks)
}

func (q *queue) runTask(task *Task) {
	defer func() {
		close(task.Err)
	}()

	for i := 0; i < q.retries; i++ {
		err := task.Fn()
		if err == nil {
			return
		}

		select {
		case task.Err <- err:
			time.Sleep(q.sleep)
		case <-time.After(q.sleep):
		}

	}
}
