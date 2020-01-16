package cache

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := assert.New(t)

	// Connection
	r, err := NewRedis("namespace")
	require.Nil(t, err)

	// Data
	data := "hello"

	// Set
	err = r.Set("key", data, 0)
	if !assert.Nil(err) {
		err, _ := err.(errors.Error)
		t.Log(err.Cause)
	}

	// Get
	v, err := r.Get("key")
	assert.Nil(err)
	if assert.NotNil(v) {
		resData, ok := v.(string)
		assert.True(ok)
		assert.Equal(data, resData)
	}

	// Delete
	err = r.Delete("key")
	assert.Nil(err)

	v, err = r.Get("key")
	assert.NotNil(err)
	assert.Nil(v)

	// Set struct
	err = r.Set("key", &struct{ value string }{"hello"}, 0)
	assert.NotNil(err)
}
