package pkg

import (
	"errors"
	"os"
)

type LocalboxConfig struct {
	Port       int    `json:"port" help:"Port to listen on" short:"p" default:"2000"`
	EngineRoot string `json:"engine_root" help:"Path where engine definitions are stored" default:"/lib/localbox/engines"`
	PoolSize   int    `json:"pool_size" help:"Total number of sandboxes that can be used concurrently" default:"10"`
	IsolateBin string `json:"isolate_bin" help:"Path to the isolate binary" default:"/usr/local/bin/isolate"`
}

type globals struct {
	EngineManager *EngineManager
	SandboxPool   *SandboxPool
	IsolateBin    string
	ShellBin      string
}

var Globals globals

func SetupLocalbox(options *LocalboxConfig) error {
	Globals = globals{
		EngineManager: &EngineManager{Index: options.EngineRoot},
		SandboxPool:   NewSandboxPool(options.PoolSize),
		IsolateBin:    options.IsolateBin,
		ShellBin:      os.Getenv("SHELL"),
	}

	if Globals.ShellBin == "" {
		return errors.New("SHELL environment variable is not set")
	}

	return nil
}
