package pkg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/pkg"
)

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

	for engineName := range engines {
		engine, err := pkg.Globals.EngineManager.Get(engineName)
		require.NoError(t, err)
		t.Run(engineName, func(t *testing.T) {
			t.Parallel()
			sandbox, err := pkg.Globals.SandboxPool.Acquire()
			require.NoError(t, err)
			t.Cleanup(func() {
				require.NoError(t, pkg.Globals.SandboxPool.Release(sandbox))
			})
			require.NoError(t, sandbox.Mount(getFiles(t, engineName)))

			// install the engine
			require.NoError(t, engine.Install())

			options := new(pkg.SandboxPhaseOptions)
			defaults.SetDefaults(options)
			options.Stdin = "localbox"

			// run the engine
			result, err := engine.Run(sandbox, options)

			require.NoError(t, err)
			require.Equal(t, "OK", result.Status, result.Message, result.Stderr)
			require.Equal(t, "localbox", result.Stdout)
		})
	}
}

func TestEngineCompileError(t *testing.T) {
	engine, err := pkg.Globals.EngineManager.Get("c")
	require.NoError(t, err)
	require.NoError(t, engine.Install())

	sandbox := pkg.Sandbox(0)
	require.NoError(t, sandbox.Init())
	t.Cleanup(func() {
		require.NoError(t, sandbox.Cleanup())
	})
	require.NoError(t, sandbox.Mount([]pkg.SandboxFile{
		{
			Name:    "main.c",
			Content: "this should throw a compile error",
		},
	}))

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	result, err := engine.Run(sandbox, options)

	require.NoError(t, err)
	require.Equal(t, "CE", result.Status)
	require.Contains(t, result.Stderr, "main.c:1:1")
}
