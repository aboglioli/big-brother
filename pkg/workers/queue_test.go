package workers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	assert := assert.New(t)

	queue := &queue{
		tasks:   make(chan *task),
		retries: 10,
		sleep:   3 * time.Millisecond,
	}
	defer func() {
		queue.Finish()
	}()
	go queue.Run()

	t1 := queue.Do(func() error {
		return errors.New("error1")
	})
	t2 := queue.Do(func() error {
		return errors.New("error2")
	})

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		c2 := 0
		for err2 := range t2.Err() {
			c2++
			assert.Equal(errors.New("error2"), err2)
		}
		assert.Equal(10, c2)
		wg.Done()
	}()

	go func() {
		c1 := 0
		for err1 := range t1.Err() {
			c1++
			assert.Equal(errors.New("error1"), err1)
		}
		assert.Equal(10, c1)
		wg.Done()
	}()

	go func() {
		d1 := <-t1.Done()
		d2 := <-t2.Done()
		assert.False(d1)
		assert.False(d2)
		wg.Done()
	}()

	wg.Wait()
}
