package internal

import (
	"errors"
	"fmt"
	"log"

	"github.com/thecomputerm/localbox/pkg"
)

var EngineManager = pkg.EngineManager{
	Index: "/workspaces/localbox/engines",
}
var SandboxPool = pkg.NewSandboxPool(10)

func init() {
	if err := SandboxPool.InitCGroup(); err != nil {
		log.Fatal(errors.Join(fmt.Errorf("couldn't init cgroup"), err))
	}
}
