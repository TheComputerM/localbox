package routes

import (
	"context"
	"runtime/debug"

	"github.com/thecomputerm/localbox/pkg"
)

type SystemInfoResponse struct {
	Body struct {
		Version       string              `json:"version" doc:"Version of localbox"`
		Configuration *pkg.LocalboxConfig `json:"configuration" doc:"Variables with which localbox was configured"`
	}
}

func GetSystemInfo(ctx context.Context, options *pkg.LocalboxConfig) (*SystemInfoResponse, error) {
	output := &SystemInfoResponse{}
	info, ok := debug.ReadBuildInfo()
	if ok {
		output.Body.Version = info.Main.Version
	}
	output.Body.Configuration = options
	return output, nil
}
