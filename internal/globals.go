package internal

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/thecomputerm/localbox/pkg"
)

var EngineManager = pkg.EngineManager{
	Index: "/workspaces/localbox/engines",
}
var SandboxPool = pkg.NewSandboxPool(10)

func init() {
	if os.Getuid() != 0 {
		log.Fatal("LocalBox must be run as root")
	}
	if err := SandboxPool.InitCGroup(); err != nil {
		log.Fatal(errors.Join(fmt.Errorf("couldn't init cgroup"), err))
	}
}
