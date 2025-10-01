package pkg_test

import (
	"errors"
	"os"

	"github.com/thecomputerm/localbox/internal"
	"github.com/thecomputerm/localbox/pkg"
)

// initializes localbox for tests
func init() {
	if err := internal.InitCGroup(); err != nil {
		panic(errors.Join(errors.New("failed to init cgroup"), err))
	}
	options := &pkg.LocalboxConfig{
		EngineRoot: os.Getenv("SERVICE_ENGINE_ROOT"),
		IsolateBin: os.Getenv("SERVICE_ISOLATE_BIN"),
		PoolSize:   1,
	}
	if err := pkg.SetupLocalbox(options); err != nil {
		panic(err)
	}
}
