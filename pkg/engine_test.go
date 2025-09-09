package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/pkg"
)

func getTestEngineManager() *pkg.EngineManager {
	return &pkg.EngineManager{
		Index: "/workspaces/localbox/engines",
	}
}

func TestGetEngine(t *testing.T) {
	manager := getTestEngineManager()

	engine, err := manager.Get("python")
	require.NoError(t, err)
	require.NotNil(t, engine)
}
