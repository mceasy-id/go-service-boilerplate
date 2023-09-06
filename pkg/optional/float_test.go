package optional_test

import (
	"encoding/json"
	"testing"

	"mceasy/service-demo/pkg/optional"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypes_NullFloat32Unmarshal(t *testing.T) {
	type Request struct {
		Name   string           `json:"name"`
		Weight optional.Float32 `json:"age"`
	}
	t.Run("invalid when key doesnt exists", func(t *testing.T) {
		input := `{
			"name": ""
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.False(t, request.Weight.IsValueSet())
	})

	t.Run("valid when null float", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":null
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Weight.IsValueSet())
		require.False(t, request.Weight.IsPresent())
	})
	t.Run("valid when 0", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":0
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Weight.IsValueSet())
		require.True(t, request.Weight.IsPresent())
		require.Equal(t, float32(0), request.Weight.MustGet())
	})
	t.Run("valid when non zero float", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":2
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Weight.IsValueSet())
		require.True(t, request.Weight.IsPresent())
		require.Equal(t, float32(2), request.Weight.MustGet())
	})
}
func TestTypes_NullFloat32Scan(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var float32Val optional.Float32
		err := float32Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, float32Val.IsValueSet())
		assert.False(t, float32Val.IsPresent())
	})
	t.Run("exists but empty", func(t *testing.T) {
		var float32Val optional.Float32
		err := float32Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, float32Val.IsValueSet())
		assert.False(t, float32Val.IsPresent())
	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var float32Val optional.Float32
		err := float32Val.Scan(float64(32.19))
		require.NoError(t, err)
		require.True(t, float32Val.IsPresent())

		assert.Equal(t, float32(32.19), float32Val.MustGet())
	})
}

func TestTypes_NullFloat32Value(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var float32Val optional.Float32
		value, err := float32Val.Value()
		require.NoError(t, err)
		assert.Empty(t, value)
		assert.Equal(t, nil, value)
	})
	t.Run("exists but empty", func(t *testing.T) {
		var float32Val optional.Float32
		float32Val.SetEmpty()
		value, err := float32Val.Value()
		require.NoError(t, err)
		assert.Empty(t, value)
		require.Nil(t, value)

	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var float32Val optional.Float32
		float32Val.Set(2.19)
		value, err := float32Val.Value()
		require.NoError(t, err)
		require.NotNil(t, value)
		_, ok := value.(float32)
		assert.True(t, ok)
		assert.Equal(t, float32(2.19), float32(value.(float32)))
	})
}
func TestTypes_NullFloat64Unmarshal(t *testing.T) {
	type Request struct {
		Name   string           `json:"name"`
		Weight optional.Float64 `json:"age"`
	}
	t.Run("invalid when key doesnt exists", func(t *testing.T) {
		input := `{
			"name": ""
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.False(t, request.Weight.IsValueSet())
	})

	t.Run("valid when null float", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":null
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Weight.IsValueSet())
		require.False(t, request.Weight.IsPresent())
	})
	t.Run("valid when 0", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":0
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Weight.IsValueSet())
		require.True(t, request.Weight.IsPresent())
		require.Equal(t, float64(0), request.Weight.MustGet())
	})
	t.Run("valid when non zero float", func(t *testing.T) {
		input := `{
			"name": "cool",
			"age":2
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Weight.IsValueSet())
		require.True(t, request.Weight.IsPresent())
		require.Equal(t, float64(2), request.Weight.MustGet())
	})
}
func TestTypes_NullFloat64Scan(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var float64Val optional.Float64
		err := float64Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, float64Val.IsValueSet())
		assert.False(t, float64Val.IsPresent())
	})
	t.Run("exists but empty", func(t *testing.T) {
		var float64Val optional.Float64
		err := float64Val.Scan(nil)
		require.NoError(t, err)
		assert.False(t, float64Val.IsValueSet())
		assert.False(t, float64Val.IsPresent())
	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var float64Val optional.Float64
		err := float64Val.Scan(float64(32.19))
		require.NoError(t, err)
		assert.True(t, float64Val.IsValueSet())
		require.True(t, float64Val.IsPresent())
		assert.Equal(t, float64(32.19), float64Val.MustGet())
	})
}

func TestTypes_NullFloat64Value(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var float32Val optional.Float64
		value, err := float32Val.Value()
		require.NoError(t, err)
		assert.Empty(t, value)
		assert.Equal(t, nil, value)
	})
	t.Run("exists but empty", func(t *testing.T) {
		var float32Val optional.Float64
		float32Val.SetEmpty()
		value, err := float32Val.Value()
		require.NoError(t, err)
		assert.Empty(t, value)
		require.Nil(t, value)

	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var float32Val optional.Float64
		float32Val.Set(2.19)
		value, err := float32Val.Value()
		require.NoError(t, err)
		require.NotNil(t, value)
		assert.Equal(t, 2.19, value)
	})
}
