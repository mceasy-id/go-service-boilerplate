package optional_test

import (
	"encoding/json"
	"testing"

	"mceasy/service-demo/pkg/optional"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypes_NullStringUnmarshal(t *testing.T) {
	type Request struct {
		Name optional.String `json:"name"`
		Age  int             `json:"age"`
	}
	t.Run("exists if value is null", func(t *testing.T) {
		input := `{
			"name": null,
			"age": 21
		}`
		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Name.IsValueSet())
		require.False(t, request.Name.IsPresent())
	})
	t.Run("exists if value is empty string", func(t *testing.T) {
		input := `{
			"name": "",
			"age": 21
		}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Name.IsValueSet())
		require.True(t, request.Name.IsPresent())
		require.Equal(t, "", request.Name.MustGet())
	})

	t.Run("exists if field is filled", func(t *testing.T) {
		input := `{
			"name": "John",
			"age": 21
			}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.True(t, request.Name.IsValueSet())
		require.NotEmpty(t, request.Name.IsPresent())
		assert.Equal(t, "John", request.Name.MustGet())
	})
	t.Run("doesnt exists if key is not exists", func(t *testing.T) {
		input := `{
					"age": 21
				}`

		var request Request

		err := json.Unmarshal([]byte(input), &request)
		require.NoError(t, err)

		require.False(t, request.Name.IsValueSet())
		require.False(t, request.Name.IsPresent())
	})
}
func TestTypes_NullStringValue(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var stringVal optional.String
		value, err := stringVal.Value()
		require.NoError(t, err)
		assert.Empty(t, value)
		assert.Equal(t, nil, value)
	})
	t.Run("exists but empty", func(t *testing.T) {
		var stringVal optional.String
		stringVal.Set("")
		value, err := stringVal.Value()
		require.NoError(t, err)
		assert.Empty(t, value)
		assert.Equal(t, "", value)

	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var stringVal optional.String
		stringVal.Set("golang the best")
		value, err := stringVal.Value()
		require.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.Equal(t, "golang the best", value)
	})

	t.Run("marshal", func(t *testing.T) {
		input := `{"name": "test"}`

		type req struct {
			Name optional.String `json:"name"`
		}

		var request req

		json.Unmarshal([]byte(input), &request)

		result, err := json.Marshal(request)
		require.NoError(t, err)
		require.NotEmpty(t, result)
	})
}
func TestTypes_NullStringScan(t *testing.T) {
	t.Run("not exists", func(t *testing.T) {
		var stringVal optional.String
		err := stringVal.Scan(nil)
		require.NoError(t, err)
		assert.False(t, stringVal.IsValueSet())
		assert.False(t, stringVal.IsPresent())
	})
	t.Run("exists but empty", func(t *testing.T) {
		var stringVal optional.String
		err := stringVal.Scan("")
		require.NoError(t, err)
		assert.True(t, stringVal.IsValueSet())
		require.True(t, stringVal.IsPresent())
		assert.Equal(t, "", stringVal.MustGet())
	})
	t.Run("exists and value not empty", func(t *testing.T) {
		var stringVal optional.String
		err := stringVal.Scan("golang the best")
		require.NoError(t, err)
		assert.True(t, stringVal.IsValueSet())
		require.True(t, stringVal.IsPresent())
		assert.Equal(t, "golang the best", stringVal.MustGet())
	})
	t.Run("exists and different data type", func(t *testing.T) {
		var stringVal optional.String
		err := stringVal.Scan(1234)
		require.NoError(t, err)
		assert.True(t, stringVal.IsValueSet())
		require.NotNil(t, stringVal.IsPresent())
		assert.Equal(t, "1234", stringVal.MustGet())
	})
}
