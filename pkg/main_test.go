package pkg_test

import (
	"errors"
	"os"
	"runtime"
	"testing"

	"github.com/thecomputerm/localbox/internal"
	"github.com/thecomputerm/localbox/pkg"
)

func TestMain(m *testing.M) {
	if err := internal.InitCGroup(); err != nil {
		panic(errors.Join(errors.New("failed to init cgroup"), err))
	}
	if err := pkg.SetOptions(&pkg.LocalboxOptions{
		EngineRoot: os.Getenv("SERVICE_ENGINE_ROOT"),
		IsolateBin: os.Getenv("SERVICE_ISOLATE_BIN"),
		PoolSize:   runtime.GOMAXPROCS(0),
	}); err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}
