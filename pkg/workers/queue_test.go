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
	c3 := 0
	t3 := queue.Do(func() error {
		t.Log(c3)
		if c3 == 2 {
			return nil
		}
		c3++
		return errors.New("error3")
	})

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		c := 0
		for err2 := range t2.Err() {
			c++
			assert.Equal(errors.New("error2"), err2)
		}
		assert.Equal(10, c)
		wg.Done()
	}()

	go func() {
		c := 0
		for err1 := range t1.Err() {
			c++
			assert.Equal(errors.New("error1"), err1)
		}
		assert.Equal(10, c)
		wg.Done()
	}()

	go func() {
		c := 0
		for err1 := range t3.Err() {
			c++
			assert.Equal(errors.New("error3"), err1)
		}
		assert.Equal(2, c)
		wg.Done()
	}()

	go func() {
		select {
		case d := <-t1.Done():
			assert.False(d)
		case d := <-t2.Done():
			assert.False(d)
		case d := <-t3.Done():
			assert.True(d)
		}
		wg.Done()
	}()

	wg.Wait()
}
