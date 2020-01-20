package workers

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueueRetries(t *testing.T) {
	getRes := func(chRes <-chan Result) Result {
		select {
		case res := <-chRes:
			return res
		case <-time.After(3 * time.Second):
			t.Error("hanging result")
		}
		return Result{}
	}

	t.Run("completed after a few retries", func(t *testing.T) {
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

		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes))
		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes))
		assert.Equal(Result{true, nil, true}, getRes(chRes))
		assert.Equal(3, cont)

		_, ok := <-chRes
		assert.False(ok)
	})

	t.Run("not completed task", func(t *testing.T) {
		assert := assert.New(t)
		q := NewQueue(4, 4, 3*time.Millisecond)
		defer q.Finish()
		q.Run()

		cont := 0
		chRes := q.Do(func() error {
			cont++
			return errors.New("err")
		})

		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes))
		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes))
		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes))
		assert.Equal(Result{true, errors.New("err"), false}, getRes(chRes))
		assert.Equal(4, cont)
	})

	t.Run("hanging task", func(t *testing.T) {
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

		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes2))
		assert.Equal(Result{false, errors.New("err"), false}, getRes(chRes2))
		assert.Equal(Result{true, nil, true}, getRes(chRes2))
	})

}
