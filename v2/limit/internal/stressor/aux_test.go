package stressor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRandom(t *testing.T) {
	const amount = 1024

	data, err := getRandom(amount)
	require.NoError(t, err)
	require.Equal(t, amount, len(data))
}
