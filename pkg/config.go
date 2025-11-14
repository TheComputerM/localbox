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
}

var instance atomic.Value

type globals struct {
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

	instance.Store(globals{
		EngineManager: &EngineManager{Root: engineRoot},
		SandboxPool:   NewSandboxPool(options.PoolSize),
		IsolateBin:    options.IsolateBin,
		ShellBin:      os.Getenv("SHELL"),
	})

	return nil
}

func Instance() globals {
	value, ok := instance.Load().(globals)
	if !ok {
		panic("globals instance not set")
	}
	return value
}
