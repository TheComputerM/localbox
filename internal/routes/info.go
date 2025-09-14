package routes

import (
	"context"

	"github.com/thecomputerm/localbox/internal"
)

type InfoResponse struct {
	Body struct {
		EngineRoot         string `json:"engine_root" doc:"Path where engine definitions are stored"`
		SandboxPoolSize    int    `json:"sandbox_pool_size" doc:"Total number of sandboxes that can be used concurrently"`
		AvailableSandboxes int    `json:"available_sandboxes" doc:"Number of sandboxes from the pool that are free at this moment"`
	}
}

func GetSystemInfo(ctx context.Context, _ *struct{}) (*InfoResponse, error) {
	output := &InfoResponse{}
	output.Body.EngineRoot = internal.EngineManager.Index
	output.Body.SandboxPoolSize = internal.SandboxPool.Capacity()
	output.Body.AvailableSandboxes = internal.SandboxPool.Available()
	return output, nil
}
