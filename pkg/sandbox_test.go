package pkg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/pkg"
)

func TestSandboxFileAccess(t *testing.T) {
	box := pkg.Sandbox(0)
	box.Init()
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	require.NoError(t, box.Mount([]pkg.SandboxFile{
		{
			Name:    "hello.txt",
			Content: "Hello World",
		},
	}))

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	result, err := box.Run(
		&pkg.SandboxPhase{
			Command:  "cat hello.txt",
			Packages: []string{"nixpkgs/nixos-25.05#busybox"},
		},
		options,
	)
	require.NoError(t, err)

	require.Equal(t, "OK", result.Status)
	require.Equal(t, 0, result.ExitCode)
	require.Equal(t, "Hello World", result.Stdout)
	require.Equal(t, "", result.Stderr)
}

func TestStdin(t *testing.T) {
	box := pkg.Sandbox(0)
	box.Init()
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	options.Stdin = "Hello World"
	result, err := box.Run(
		&pkg.SandboxPhase{
			Command:  "tee output.txt",
			Packages: []string{"nixpkgs/nixos-25.05#busybox"},
		},
		options,
	)
	require.NoError(t, err)
	require.Equal(t, "OK", result.Status)
	require.Equal(t, 0, result.ExitCode)
	require.Equal(t, "", result.Stderr)

	outputFile := filepath.Join(box.BoxPath(), "output.txt")
	require.FileExists(t, outputFile)
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "Hello World", string(content))
}

func TestWallTime(t *testing.T) {
	box := pkg.Sandbox(0)
	box.Init()
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	result, err := box.Run(
		&pkg.SandboxPhase{
			Command:  "sleep 2",
			Packages: []string{"nixpkgs/nixos-25.05#busybox"},
		},
		options,
	)
	require.NoError(t, err)

	require.Equal(t, "OK", result.Status)
	require.Greater(t, result.WallTime, 2000)
}
