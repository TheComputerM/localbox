package internal

import "github.com/thecomputerm/localbox/pkg"

var EngineManager = pkg.EngineManager{
	Index: "/workspaces/localbox/engines",
}
var SandboxPool = pkg.NewSandboxPool(10)

// func init() {
// 	if err := SandboxPool.InitCGroup(); err != nil {
// 		panic(err)
// 	}
// }
