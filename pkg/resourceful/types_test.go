package resourceful

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// this is only for dev, if in the future you want to refactor and this test is bugging you,
// you can just delete it, its just a test
func TestResource_sanitizeFilterValueByType(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		satizedFilterValue, ok := sanitizeFilterValueByType("(\"halo' semuanya\" \"coba'\" 'dulu\"' '\"ayam geprek' kucing \"'bededeh'\" '\"ayamku\"' mantap)", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"halo' semuanya", "coba'", "dulu\"", "\"ayam geprek", "kucing", "'bededeh'", "\"ayamku\"", "mantap"}, satizedFilterValue)

		satizedFilterValue, ok = sanitizeFilterValueByType("(\"perabotan dan alat rumah tangga\" hobi 'elektronik' 'makanan dan minuman')", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"perabotan dan alat rumah tangga", "hobi", "elektronik", "makanan dan minuman"}, satizedFilterValue)

		satizedFilterValue, ok = sanitizeFilterValueByType("(makanan)", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"makanan"}, satizedFilterValue)

		satizedFilterValue, ok = sanitizeFilterValueByType("makanan minuman", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"makanan", "minuman"}, satizedFilterValue)

		satizedFilterValue, ok = sanitizeFilterValueByType("minuman \"mak\"anan\"", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"minuman", "mak\"anan"}, satizedFilterValue)

		satizedFilterValue, ok = sanitizeFilterValueByType("minuman 'mak'anan'", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"minuman", "mak'anan"}, satizedFilterValue)

		// not supported
		satizedFilterValue, ok = sanitizeFilterValueByType("('VELI'S x Bagenciala T-Shirt' \"Shoes \"Air\" Jordan\")", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"VELI'S x Bagenciala T-Shirt", "Shoes \"Air", "Jordan\""}, satizedFilterValue)
	})

	t.Run("invalid", func(t *testing.T) {
		//TODO: handle validation on one quoted string
		satizedFilterValue, ok := sanitizeFilterValueByType("minuman 'makanan", STRING)
		require.True(t, ok)
		require.Equal(t, []string{"minuman"}, satizedFilterValue)
	})

}
