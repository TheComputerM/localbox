package pkg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/pkg"
)

func TestSandboxMount(t *testing.T) {
	box := pkg.Sandbox(0)
	require.NoError(t, box.Init())
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	t.Run("escape", func(t *testing.T) {
		require.Error(t, box.Mount([]pkg.SandboxFile{
			{
				Name:    "../escape.txt",
				Content: "This should not work",
			},
		}))
	})

	t.Run("simple", func(t *testing.T) {
		require.NoError(t, box.Mount([]pkg.SandboxFile{
			{
				Name:    "file1.txt",
				Content: "1",
			},
			{
				Name:    "./file2.txt",
				Content: "2",
			},
		}))
		file1 := filepath.Join(box.BoxPath(), "file1.txt")
		require.FileExists(t, file1)
		content, err := os.ReadFile(file1)
		require.NoError(t, err)
		require.Equal(t, "1", string(content))

		file2 := filepath.Join(box.BoxPath(), "file2.txt")
		require.FileExists(t, file2)
		content, err = os.ReadFile(file2)
		require.NoError(t, err)
		require.Equal(t, "2", string(content))
	})

	t.Run("nested", func(t *testing.T) {
		require.NoError(t, box.Mount([]pkg.SandboxFile{
			{
				Name:    "nested/file.txt",
				Content: "nested",
			},
		}))
		file := filepath.Join(box.BoxPath(), "nested", "file.txt")
		require.FileExists(t, file)
		content, err := os.ReadFile(file)
		require.NoError(t, err)
		require.Equal(t, "nested", string(content))
	})
}

func TestSandboxFileAccess(t *testing.T) {
	box := pkg.Sandbox(0)
	require.NoError(t, box.Init())
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
			SandboxCommand: pkg.SandboxCommand{
				Command:  "cat hello.txt",
				Packages: []string{"nixpkgs/nixos-25.05#busybox"},
			},
		},
		options,
	)
	require.NoError(t, err)

	require.Equal(t, "OK", result.Status)
	require.Equal(t, 0, result.ExitCode)
	require.Equal(t, "Hello World", result.Stdout)
}

func TestSandbox_Stdin(t *testing.T) {
	box := pkg.Sandbox(0)
	require.NoError(t, box.Init())
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	options.Stdin = "Hello World"
	result, err := box.Run(
		&pkg.SandboxPhase{
			SandboxCommand: pkg.SandboxCommand{
				Command:  "tee output.txt",
				Packages: []string{"nixpkgs/nixos-25.05#busybox"},
			},
		},
		options,
	)
	require.NoError(t, err)
	require.Equal(t, "OK", result.Status)
	require.Equal(t, 0, result.ExitCode)

	outputFile := filepath.Join(box.BoxPath(), "output.txt")
	require.FileExists(t, outputFile)
	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, "Hello World", string(content))
}

func TestSandbox_WallTimeLimit(t *testing.T) {
	box := pkg.Sandbox(0)
	require.NoError(t, box.Init())
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	options.WallTimeLimit = 1000
	result, err := box.Run(
		&pkg.SandboxPhase{
			SandboxCommand: pkg.SandboxCommand{
				Command:  "sleep 2",
				Packages: []string{"nixpkgs/nixos-25.05#busybox"},
			},
		},
		options,
	)
	require.NoError(t, err)

	require.Equal(t, "TO", result.Status)
}

func TestSandbox_FileSizeLimit(t *testing.T) {
	box := pkg.Sandbox(0)
	require.NoError(t, box.Init())
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	options := new(pkg.SandboxPhaseOptions)
	defaults.SetDefaults(options)
	options.FileSizeLimit = 1

	result, err := box.Run(
		&pkg.SandboxPhase{
			SandboxCommand: pkg.SandboxCommand{
				Command:  "truncate -s 2K filename.txt",
				Packages: []string{"nixpkgs/nixos-25.05#busybox"},
			},
		},
		options,
	)
	require.NoError(t, err)

	require.Equal(t, "SG", result.Status)
	require.Equal(t, 25, result.ExitCode)
}

func TestSandbox_UnsafeRun(t *testing.T) {
	box := pkg.Sandbox(0)
	require.NoError(t, box.Init())
	t.Cleanup(func() {
		require.NoError(t, box.Cleanup())
	})

	require.NoError(t, box.Mount([]pkg.SandboxFile{
		{
			Name:    "hello.txt",
			Content: "The file should be renamed to modified.txt",
		},
	}))

	_, _, err := box.UnsafeRun(&pkg.SandboxCommand{
		Command:  "mv hello.txt modified.txt",
		Packages: []string{"nixpkgs/nixos-25.05#busybox"},
	})
	require.NoError(t, err)

	require.FileExists(t, filepath.Join(box.BoxPath(), "modified.txt"))
	require.NoFileExists(t, filepath.Join(box.BoxPath(), "hello.txt"))
}
