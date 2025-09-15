package routes

import (
	"context"

	"github.com/thecomputerm/localbox/internal"
)

type SystemInfoResponse struct {
	Body struct {
		Configuration *internal.LocalboxConfig `json:"configuration" doc:"Variables with which localbox was configured"`
	}
}

func GetSystemInfo(ctx context.Context, options *internal.LocalboxConfig) (*SystemInfoResponse, error) {
	output := &SystemInfoResponse{}
	output.Body.Configuration = options
	return output, nil
}
