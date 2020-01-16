package cache

import (
	"encoding/json"
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Data struct {
	Value string `json:"value"`
}

func TestRedis(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := assert.New(t)

	// Connection
	r, err := NewRedis("namespace")
	require.Nil(t, err)

	t.Run("string", func(t *testing.T) {
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
	})

	t.Run("save struct", func(t *testing.T) {
		data := &Data{"hello"}

		// Set struct
		err = r.Set("key", data, 0)
		assert.NotNil(err)

		// Set marshalled struct
		b, err := json.Marshal(data)
		require.Nil(t, err)
		require.NotEmpty(t, b)

		err = r.Set("data123", b, 0)
		assert.Nil(err)

		v, err := r.Get("data123")
		assert.Nil(err)
		if assert.NotNil(v) {
			s, ok := v.(string)
			assert.True(ok)

			d := &Data{}
			err := json.Unmarshal([]byte(s), d)
			assert.Nil(err)

			assert.Equal(data, d)
		}
	})
}
