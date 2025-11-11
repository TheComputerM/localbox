package pkg

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
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

func init() {
	if err := SetupLocalbox(&LocalboxConfig{
		EngineRoot: os.Getenv("SERVICE_ENGINE_ROOT"),
		IsolateBin: os.Getenv("SERVICE_ISOLATE_BIN"),
		PoolSize:   runtime.GOMAXPROCS(0),
	}); err != nil {
		panic(err)
	}

}

func SetupLocalbox(options *LocalboxConfig) error {
	engineRoot, err := filepath.Abs(options.EngineRoot)
	if err != nil {
		return err
	}

	Globals = globals{
		EngineManager: &EngineManager{Root: engineRoot},
		SandboxPool:   NewSandboxPool(options.PoolSize),
		IsolateBin:    options.IsolateBin,
		ShellBin:      os.Getenv("SHELL"),
	}

	if Globals.ShellBin == "" {
		return errors.New("SHELL environment variable is not set")
	}

	return nil
}
