package routes

import (
	"context"

	"github.com/thecomputerm/localbox/pkg"
)

type ExecutePhase struct {
	pkg.SandboxPhase
	pkg.SandboxPhaseOptions
}

type ExecuteRequest struct {
	Body struct {
		Prepare []pkg.SandboxCommand `json:"prepare" doc:"Commands to run to prepare the sandbox before execution" required:"false"`
		Phases  []ExecutePhase       `json:"phases" doc:"Execution phases to run sequentially in the sandbox"`
		Files   []pkg.SandboxFile    `json:"files" doc:"Files to mount in the sandbox before execution"`
	}
}

type ExecuteResponse struct {
	Body []*pkg.SandboxPhaseResults
}

func Execute(ctx context.Context, input *ExecuteRequest) (*ExecuteResponse, error) {
	pool := pkg.Instance().SandboxPool
	sandbox, err := pool.Acquire()
	if err != nil {
		return nil, err
	}
	defer pool.Release(sandbox)
	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}

	for _, prep := range input.Body.Prepare {
		if result, ok := sandbox.UnsafeRun(&prep); !ok {
			return &ExecuteResponse{Body: []*pkg.SandboxPhaseResults{result}}, nil
		}
	}

	results := make([]*pkg.SandboxPhaseResults, len(input.Body.Phases))
	for i, phase := range input.Body.Phases {
		result, err := sandbox.Run(&phase.SandboxPhase, &phase.SandboxPhaseOptions)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return &ExecuteResponse{Body: results}, nil
}
