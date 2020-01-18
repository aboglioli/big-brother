package workers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	assert := assert.New(t)

	queue := &queue{
		tasks:   make(chan *Task),
		retries: 3,
		sleep:   10 * time.Millisecond,
	}
	defer func() {
		queue.Finish()
	}()
	go queue.Run()

	chErr1 := queue.Do(func() error {
		return errors.New("error1")
	})
	chErr2 := queue.Do(func() error {
		return errors.New("error2")
	})

	c2 := 0
	for err2 := range chErr2 {
		c2++
		assert.Equal(errors.New("error2"), err2)
	}
	assert.Equal(3, c2)

	c1 := 0
	for err1 := range chErr1 {
		c1++
		assert.Equal(errors.New("error1"), err1)
	}
	assert.Equal(0, c1)
}
