package cache

import (
	"encoding/json"
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type data struct {
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
		mD := &data{"hello"}

		// Set struct
		err = r.Set("key", mD, 0)
		assert.NotNil(err)

		// Set marshalled struct
		b, err := json.Marshal(mD)
		require.Nil(t, err)
		require.NotEmpty(t, b)

		err = r.Set("data123", b, 0)
		assert.Nil(err)

		v, err := r.Get("data123")
		assert.Nil(err)
		if assert.NotNil(v) {
			s, ok := v.(string)
			assert.True(ok)

			d := &data{}
			err := json.Unmarshal([]byte(s), d)
			assert.Nil(err)

			assert.Equal(mD, d)
		}

		err = r.Delete("data123")
		assert.Nil(err)
	})
}
