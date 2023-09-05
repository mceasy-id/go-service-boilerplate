package optional_test

import (
	"encoding/json"
	"testing"

	"mceasy/service-demo/pkg/optional"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypes_NullBoolUnmarshal(t *testing.T) {
	type Request struct {
		Name     string        `json:"name"`
		IsLocked optional.Bool `json:"is_locked"`
	}
	t.Run("invalid when key doesnt exists", func(t *testing.T) {
		input := `{
			"name":"uncle bob"
		}`
		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		assert.False(t, request.IsLocked.IsValueSet())
		require.True(t, request.IsLocked.IsPresent())
		assert.False(t, request.IsLocked.MustGet())
	})
	t.Run("valid when value is null", func(t *testing.T) {
		input := `{
			"name":"uncle bob",
			"is_locked":null
		}`
		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		assert.True(t, request.IsLocked.IsValueSet())
		assert.True(t, request.IsLocked.IsPresent())
		assert.False(t, request.IsLocked.MustGet())
	})
	t.Run("valid when value is false", func(t *testing.T) {
		input := `{
			"name":"uncle bob",
			"is_locked":false
		}`
		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		assert.True(t, request.IsLocked.IsValueSet())
		require.True(t, request.IsLocked.IsPresent())
		assert.False(t, request.IsLocked.MustGet())
	})
	t.Run("valid when value is true", func(t *testing.T) {
		input := `{
			"name":"uncle bob",
			"is_locked":true
		}`
		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		assert.True(t, request.IsLocked.IsValueSet())
		require.True(t, request.IsLocked.IsPresent())
		assert.True(t, request.IsLocked.MustGet())
	})
}
func TestTypes_NullBoolScan(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var boolVal optional.Bool
		err := boolVal.Scan(nil)
		require.NoError(t, err)
		assert.False(t, boolVal.IsValueSet())
		require.True(t, boolVal.IsPresent())
		assert.False(t, boolVal.MustGet())
	})
	t.Run("exists but empty", func(t *testing.T) {
		var boolVal optional.Bool
		err := boolVal.Scan(false)
		require.NoError(t, err)
		assert.True(t, boolVal.IsValueSet())
		require.True(t, boolVal.IsPresent())
		assert.False(t, boolVal.MustGet())
	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var boolVal optional.Bool
		err := boolVal.Scan(true)
		require.NoError(t, err)
		assert.True(t, boolVal.IsValueSet())
		require.True(t, boolVal.IsPresent())
		assert.True(t, boolVal.MustGet())
	})
}

func TestTypes_NullBoolValue(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var boolVal optional.Bool
		value, err := boolVal.Value()
		require.NoError(t, err)
		assert.Equal(t, nil, value)
	})
	t.Run("exists but empty", func(t *testing.T) {
		var boolVal optional.Bool
		boolVal.SetEmpty()
		value, err := boolVal.Value()
		require.NoError(t, err)
		assert.Equal(t, false, value.(bool))

	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var boolVal optional.Bool
		boolVal.Set(true)
		value, err := boolVal.Value()
		require.NoError(t, err)
		assert.Equal(t, true, value.(bool))
	})
}
