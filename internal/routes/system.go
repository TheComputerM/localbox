package routes

import (
	"context"

	"github.com/thecomputerm/localbox/pkg"
)

type SystemInfoResponse struct {
	Body struct {
		Configuration *pkg.LocalboxConfig `json:"configuration" doc:"Variables with which localbox was configured"`
	}
}

func GetSystemInfo(ctx context.Context, options *pkg.LocalboxConfig) (*SystemInfoResponse, error) {
	output := &SystemInfoResponse{}
	output.Body.Configuration = options
	return output, nil
}
