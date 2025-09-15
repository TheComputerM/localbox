package internal

import (
	"github.com/thecomputerm/localbox/pkg"
)

type LocalboxConfig struct {
	Port       int    `json:"port" help:"Port to listen on" short:"p" default:"2000"`
	EngineRoot string `json:"engine_root" help:"Path where engine definitions are stored" default:"/lib/localbox/engines"`
	PoolSize   int    `json:"pool_size" help:"Total number of sandboxes that can be used concurrently" default:"10"`
}

var EngineManager *pkg.EngineManager
var SandboxPool *pkg.SandboxPool

func SetupLocalbox(options *LocalboxConfig) error {
	EngineManager = &pkg.EngineManager{
		Index: options.EngineRoot,
	}
	SandboxPool = pkg.NewSandboxPool(options.PoolSize)
	return nil
}
