package workers

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getRes(t *testing.T, chRes <-chan Result) Result {
	select {
	case res := <-chRes:
		return res
	case <-time.After(3 * time.Second):
		t.Error("hanging result")
	}
	return Result{}
}

func TestTaskCompletedAfterAFewRetries(t *testing.T) {
	assert := assert.New(t)
	q := NewQueue(4, 4, 3*time.Millisecond).(*queue)
	defer q.Finish()
	q.Run()

	cont := 0
	chRes := q.Do(func() error {
		cont++
		if cont == 3 {
			return nil
		}
		return errors.New("err")
	})
	assert.Len(q.tasks, 1)

	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes))
	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes))
	assert.Equal(Result{true, nil}, getRes(t, chRes))
	assert.Equal(3, cont)

	_, ok := <-chRes
	assert.False(ok)

}

func TestNotCompletedTask(t *testing.T) {
	assert := assert.New(t)
	q := NewQueue(4, 4, 3*time.Millisecond)
	defer q.Finish()
	q.Run()

	cont := 0
	chRes := q.Do(func() error {
		cont++
		return errors.New("err")
	})

	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes))
	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes))
	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes))
	assert.Equal(Result{true, errors.New("err")}, getRes(t, chRes))
	assert.Equal(4, cont)

}

func TestHangingTasks(t *testing.T) {
	assert := assert.New(t)
	q := NewQueue(2, 4, time.Millisecond).(*queue)
	defer q.Finish()
	q.Run()

	cont := 0
	chRes1 := q.Do(func() error {
		cont++
		if cont == 2 {
			return nil
		}
		return errors.New("err")
	})
	assert.Len(q.tasks, 1)
	_, ok := <-chRes1
	assert.True(ok)

	q.Do(func() error {
		return errors.New("err")
	})
	assert.Len(q.tasks, 2)

	q.Do(func() error {
		return nil
	})
	assert.Len(q.tasks, 2)
	_, ok = <-chRes1
	assert.False(ok)

	cont = 0
	chRes2 := q.Do(func() error {
		cont++
		if cont == 3 {
			return nil
		}
		return errors.New("err")
	})
	assert.Len(q.tasks, 2)
	_, ok = <-chRes1
	assert.False(ok)

	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes2))
	assert.Equal(Result{false, errors.New("err")}, getRes(t, chRes2))
	assert.Equal(Result{true, nil}, getRes(t, chRes2))
}

func TestInfiniteRetries(t *testing.T) {
	assert := assert.New(t)
	q := NewQueue(4, 999, 10*time.Nanosecond).(*queue)
	q.Run()

	for i := 0; i < 10; i++ {
		cont := 0
		chRes := q.Do(func() error {
			cont++
			if cont == 100 {
				return nil
			}
			return errors.New("err")
		})
		res := getRes(t, chRes)
		firstRes := res
		var lastRes Result
		for res := range chRes {
			lastRes = res
		}

		assert.Equal(100, cont)
		assert.Equal(Result{false, errors.New("err")}, firstRes)
		assert.Equal(Result{true, nil}, lastRes)

		switch i {
		case 0:
			assert.Len(q.tasks, 1)
		case 1:
			assert.Len(q.tasks, 2)
		case 2:
			assert.Len(q.tasks, 3)
		default:
			assert.Len(q.tasks, 4)
		}
	}

}

func TestExecutionOrder(t *testing.T) {
	assert := assert.New(t)
	q := NewQueue(4, 3, 0).(*queue)
	q.Run()

	var order []string
	ch := q.Do(func() error {
		time.Sleep(10 * time.Millisecond)
		order = append(order, "t1")
		return errors.New("err")
	})
	q.Do(func() error {
		order = append(order, "t2")
		return errors.New("err")
	})

	for range ch {
	}

	assert.Equal([]string{"t2", "t2", "t2", "t1", "t1", "t1"}, order)
}
