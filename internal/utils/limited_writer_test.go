package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thecomputerm/localbox/internal/utils"
)

func TestLimitedWriter(t *testing.T) {
	t.Run("Write within capacity", func(t *testing.T) {
		lw := utils.NewLimitedWriter(10)
		n, err := lw.Write([]byte("hello"))
		require.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "hello", lw.String())
	})

	t.Run("Write over capacity", func(t *testing.T) {
		lw := utils.NewLimitedWriter(5)
		n, err := lw.Write([]byte("helloworld"))
		require.ErrorIs(t, err, utils.ErrWriteLimitExceeded)
		assert.Equal(t, 5, n)
		assert.Equal(t, "hello", lw.String())
	})

	t.Run("Write with zero capacity", func(t *testing.T) {
		lw := utils.NewLimitedWriter(0)
		n, err := lw.Write([]byte("hello"))
		require.ErrorIs(t, err, utils.ErrWriteLimitExceeded)
		assert.Equal(t, 0, n)
		assert.Equal(t, "", lw.String())
	})

	t.Run("Multiple writes within capacity", func(t *testing.T) {
		lw := utils.NewLimitedWriter(11)
		n, err := lw.Write([]byte("hello"))
		require.NoError(t, err)
		assert.Equal(t, 5, n)

		n, err = lw.Write([]byte(" "))
		require.NoError(t, err)
		assert.Equal(t, 1, n)

		n, err = lw.Write([]byte("world"))
		require.NoError(t, err)
		assert.Equal(t, 5, n)

		assert.Equal(t, "hello world", lw.String())
	})

	t.Run("Multiple writes exceeding capacity", func(t *testing.T) {
		lw := utils.NewLimitedWriter(8)
		_, err := lw.Write([]byte("hello"))
		require.NoError(t, err)

		n, err := lw.Write([]byte("world"))
		require.ErrorIs(t, err, utils.ErrWriteLimitExceeded)
		assert.Equal(t, 3, n)
		assert.Equal(t, "hellowor", lw.String())
	})

	t.Run("String method with whitespace", func(t *testing.T) {
		lw := utils.NewLimitedWriter(20)
		_, err := lw.Write([]byte("  hello world  "))
		require.NoError(t, err)
		assert.Equal(t, "hello world", lw.String())
	})
}
