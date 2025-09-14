package internal

import (
	"github.com/thecomputerm/localbox/pkg"
)

type LocalboxConfig struct {
	Port       int    `help:"Port to listen on" short:"p" default:"2000"`
	EngineRoot string `help:"Path where engine definitions are stored" default:"/workspaces/localbox/engines"`
	PoolSize   int    `help:"Total number of sandboxes that can be used concurrently" default:"10"`
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
