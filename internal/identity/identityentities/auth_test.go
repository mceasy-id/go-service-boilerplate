package identityentities_test

import (
	"testing"

	"mceasy/service-demo/internal/identity/identityentities"

	"github.com/stretchr/testify/require"
)

func TestIdentitiyAuth_Validate(t *testing.T) {
	t.Run("error if empty", func(t *testing.T) {
		err := identityentities.Credential{}.Validate()
		require.Error(t, err)
	})
	t.Run("ok if not empty", func(t *testing.T) {
		err := identityentities.Credential{
			UserName:  "calvary",
			UserId:    1,
			CompanyId: 392,
		}.Validate()
		require.NoError(t, err)
	})
}
