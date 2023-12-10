package generator

import (
	"github.com/romandnk/shortener/internal/constant"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGen_Random(t *testing.T) {
	gen := Gen{Length: constant.AliasLength}

	uniqueStrings := make(map[string]struct{})

	for i := 0; i < 10_000_000; i++ {
		randomString, err := gen.Random()
		require.NoError(t, err)

		_, exists := uniqueStrings[randomString]
		require.False(t, exists)

		uniqueStrings[randomString] = struct{}{}
	}
}
