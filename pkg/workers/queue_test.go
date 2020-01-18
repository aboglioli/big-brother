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
		tasks:      make(chan *Task),
		retries:    3,
		waitForErr: 2 * time.Second,
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

	err2 := <-chErr2
	assert.Equal(errors.New("error2"), err2)

	err1 := <-chErr1
	assert.Equal(errors.New("error1"), err1)

}
