package pkg_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/pkg"
)

func TestSandboxPool(t *testing.T) {
	pool := pkg.Instance().SandboxPool

	t.Run("one", func(t *testing.T) {
		sandbox, err := pool.Acquire()
		require.NoError(t, err)
		require.Equal(t, pool.Capacity()-1, pool.Available())
		require.NoError(t, pool.Release(sandbox))
		require.Equal(t, pool.Capacity(), pool.Available())
	})

	t.Run("blocking", func(t *testing.T) {
		sandboxes := make([]pkg.Sandbox, pool.Capacity())
		for i := 0; i < pool.Capacity(); i++ {
			sandbox, err := pool.Acquire()
			require.NoError(t, err)
			sandboxes[i] = sandbox
		}

		require.Equal(t, 0, pool.Available())
		acquired := false
		go func() {
			sandbox, err := pool.Acquire()
			require.NoError(t, err)
			sandboxes[0] = sandbox
			acquired = true
		}()

		require.False(t, acquired, "Should not have acquired a sandbox yet")
		require.NoError(t, pool.Release(sandboxes[0]))
		require.Eventually(
			t,
			func() bool { return acquired },
			time.Second,
			100*time.Millisecond,
			"Should have acquired a sandbox after release",
		)

		for i := 0; i < pool.Capacity(); i++ {
			require.NoError(t, pool.Release(sandboxes[i]))
		}
		require.Equal(t, pool.Capacity(), pool.Available())
	})
}
