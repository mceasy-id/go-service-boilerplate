package optional_test

import (
	"encoding/json"
	"testing"

	"mceasy/service-demo/pkg/optional"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOptional(t *testing.T) {
	t.Run("marshalJSON", func(t *testing.T) {
		var myStruct struct {
			ID optional.Option[uuid.UUID] `json:"id"`
		}

		err := json.Unmarshal([]byte(`
		{
			"id":null
		}
		`), &myStruct)

		require.NoError(t, err)

		require.True(t, myStruct.ID.IsValueSet())
		require.False(t, myStruct.ID.IsPresent())
	})
}
