package routes

import (
	"context"
	"os/exec"

	"github.com/thecomputerm/localbox/pkg"
)

type SystemInfoResponse struct {
	Body struct {
		AvailableSandboxes int                 `json:"available_sandboxes" example:"10" doc:"Number of sandboxes currently available"`
		Configuration      *pkg.LocalboxConfig `json:"configuration" doc:"Variables with which localbox was configured"`
	}
}

func GetSystemInfo(ctx context.Context, options *pkg.LocalboxConfig) (*SystemInfoResponse, error) {
	output := &SystemInfoResponse{}
	output.Body.AvailableSandboxes = pkg.Globals.SandboxPool.Available()
	output.Body.Configuration = options
	return output, nil
}

func RunSystemGC(ctx context.Context, _ *struct{}) (*struct{}, error) {
	cmd := exec.Command("nix", "store", "gc")
	return nil, cmd.Run()
}
