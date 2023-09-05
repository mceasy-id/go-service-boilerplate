package optional_test

import (
	"encoding/json"
	"testing"

	"mceasy/service-demo/pkg/optional"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypes_NullInt32Unmarshal(t *testing.T) {
	type Request struct {
		Name string         `json:"name"`
		Age  optional.Int32 `json:"age"`
	}
	t.Run("invalid when key doesnt exists", func(t *testing.T) {
		input := `{
			"name": ""
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.False(t, request.Age.IsValueSet())
	})

	t.Run("valid when null int", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":null
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Age.IsValueSet())
		require.False(t, request.Age.IsPresent())
	})
	t.Run("valid when 0", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":0
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Age.IsValueSet())
		require.Equal(t, int32(0), request.Age.MustGet())
	})
	t.Run("valid when non zero int", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":2
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Age.IsValueSet())
		require.Equal(t, int32(2), request.Age.MustGet())
	})
}
func TestTypes_NullInt32Scan(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var int32Val optional.Int32
		err := int32Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, int32Val.IsValueSet())
		assert.False(t, int32Val.IsPresent())
	})
	t.Run("exists but empty", func(t *testing.T) {
		var Int32Val optional.Int32
		err := Int32Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, Int32Val.IsPresent())
		assert.False(t, Int32Val.IsPresent())
	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var Int32Val optional.Int32
		err := Int32Val.Scan(23)
		require.NoError(t, err)
		assert.True(t, Int32Val.IsValueSet())
		require.True(t, Int32Val.IsPresent())
		assert.Equal(t, int32(23), Int32Val.MustGet())
	})
}

func TestTypes_NullInt32Value(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var Int32Val optional.Int32
		value, err := Int32Val.Value()
		require.NoError(t, err)
		assert.Nil(t, value)
	})
	t.Run("exists but empty", func(t *testing.T) {
		var Int32Val optional.Int32
		Int32Val.SetEmpty()
		value, err := Int32Val.Value()
		require.NoError(t, err)
		require.Nil(t, value)

	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var Int32Val optional.Int32
		Int32Val.Set(23)
		value, err := Int32Val.Value()
		require.NoError(t, err)
		require.NotNil(t, value)
		assert.Equal(t, int64(23), value)
	})
}
func TestTypes_NullInt64Unmarshal(t *testing.T) {
	type Request struct {
		Name string         `json:"name"`
		Age  optional.Int64 `json:"age"`
	}
	t.Run("invalid when key doesnt exists", func(t *testing.T) {
		input := `{
			"name": ""
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.False(t, request.Age.IsValueSet())
	})

	t.Run("valid when null int", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":null
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Age.IsValueSet())
		require.False(t, request.Age.IsPresent())
	})
	t.Run("valid when 0", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":0
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Age.IsPresent())
		require.Equal(t, int64(0), request.Age.MustGet())
	})
	t.Run("valid when non zero int", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":2
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Age.IsPresent())
		require.Equal(t, int64(2), request.Age.MustGet())
	})
}
func TestTypes_NullInt64Scan(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var int64Val optional.Int64
		err := int64Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, int64Val.IsValueSet())
	})
	t.Run("exists but empty", func(t *testing.T) {
		var int64Val optional.Int64
		err := int64Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, int64Val.IsValueSet())
		assert.False(t, int64Val.IsPresent())
	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var int64Val optional.Int64
		err := int64Val.Scan(23)
		require.NoError(t, err)
		assert.True(t, int64Val.IsValueSet())
		assert.True(t, int64Val.IsPresent())
		assert.Equal(t, int64(23), int64Val.MustGet())
	})
}

func TestTypes_NullInt64Value(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var int64Val optional.Int64
		value, err := int64Val.Value()
		require.NoError(t, err)
		assert.Nil(t, value)
	})
	t.Run("exists but empty", func(t *testing.T) {
		var int64Val optional.Int64
		int64Val.SetEmpty()
		value, err := int64Val.Value()
		require.NoError(t, err)
		require.Nil(t, value)

	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var int64Val optional.Int64
		int64Val.Set(23)
		value, err := int64Val.Value()
		require.NoError(t, err)
		require.NotNil(t, value)
		assert.Equal(t, int64(23), value)
	})
}
