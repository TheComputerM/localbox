package pkg_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/internal"
	"github.com/thecomputerm/localbox/pkg"
)

func init() {
	if err := internal.InitCGroup(); err != nil {
		panic(errors.Join(errors.New("failed to init cgroup"), err))
	}
	options := &pkg.LocalboxConfig{
		EngineRoot: os.Getenv("SERVICE_ENGINE_ROOT"),
		PoolSize:   1,
		IsolateBin: os.Getenv("SERVICE_ISOLATE_BIN"),
	}
	if err := pkg.SetupLocalbox(options); err != nil {
		panic(err)
	}
}

func TestGetEngine(t *testing.T) {
	engine, err := pkg.Globals.EngineManager.Get("python")
	require.NoError(t, err)
	require.NotNil(t, engine)
}

func TestEngines(t *testing.T) {
	testdataDir := "../test/engines"
	files, err := os.ReadDir(testdataDir)
	require.NoError(t, err)
	require.Greater(t, len(files), 0)

	// gets the files to mount for the given engine
	getFiles := func(t *testing.T, engine string) []pkg.SandboxFile {
		entries, err := os.ReadDir(filepath.Join(testdataDir, engine))
		require.NoError(t, err)
		files := make([]pkg.SandboxFile, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			content, err := os.ReadFile(filepath.Join(testdataDir, engine, name))
			require.NoError(t, err)
			files = append(files, pkg.SandboxFile{
				Name:    name,
				Content: string(content),
			})
		}
		return files
	}

	engines, err := pkg.Globals.EngineManager.List()
	require.NoError(t, err)
	require.Greater(t, len(engines), 0)
	sandbox := pkg.Sandbox(0)

	for engineName := range engines {
		engine, err := pkg.Globals.EngineManager.Get(engineName)
		require.NoError(t, err)
		t.Run(engineName, func(t *testing.T) {
			require.NoError(t, sandbox.Init())
			t.Cleanup(func() {
				require.NoError(t, sandbox.Cleanup())
			})
			require.NoError(t, sandbox.Mount(getFiles(t, engineName)))
			opts := defaultSandboxPhaseOptions()
			opts.Stdin = "localbox"
			execute, err := engine.Run(sandbox, opts)
			require.NoError(t, err)
			require.Equal(t, "OK", execute.Status)
			require.Equal(t, "localbox", execute.Stdout)
		})
	}
}
