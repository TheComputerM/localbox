package pkg

import (
	"os"
	"path/filepath"
	"sync/atomic"
)

type LocalboxOptions struct {
	EngineRoot string `json:"engine_root" help:"Path where engine definitions are stored" default:"/lib/localbox/engines"`
	PoolSize   int    `json:"pool_size" help:"Total number of sandboxes that can be used concurrently" default:"10"`
	IsolateBin string `json:"isolate_bin" help:"Path to the isolate binary" default:"/usr/local/bin/isolate"`
	ShellBin   string `json:"shell_bin" help:"Path to the shell binary used in sandboxes, defaults to the SHELL environment variable"`
}

var store atomic.Value

type instance struct {
	EngineManager *EngineManager
	SandboxPool   *SandboxPool
	IsolateBin    string
	ShellBin      string
}

func SetOptions(options *LocalboxOptions) error {
	engineRoot, err := filepath.Abs(options.EngineRoot)
	if err != nil {
		return err
	}

	config := instance{
		EngineManager: &EngineManager{Root: engineRoot},
		SandboxPool:   NewSandboxPool(options.PoolSize),
		IsolateBin:    options.IsolateBin,
		ShellBin:      options.ShellBin,
	}

	if config.ShellBin == "" {
		config.ShellBin = os.Getenv("SHELL")
	}

	store.Store(config)

	return nil
}

func Instance() instance {
	value, ok := store.Load().(instance)
	if !ok {
		panic("globals instance not set")
	}
	return value
}
