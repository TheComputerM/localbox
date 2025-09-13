package routes

import (
	"context"

	"github.com/thecomputerm/localbox/internal"
	"github.com/thecomputerm/localbox/pkg"
)

type ExecutePhase struct {
	pkg.SandboxPhase
	pkg.SandboxPhaseOptions
}

type ExecuteRequest struct {
	Body struct {
		Prepare []pkg.SandboxPrepare `json:"prepare,omitempty" doc:"Preparation steps to run before execution"`
		Phases  []ExecutePhase       `json:"phases" doc:"Execution phases to run sequentially in the sandbox"`
		Files   []pkg.SandboxFile    `json:"files" doc:"Files to mount in the sandbox before execution"`
	}
}

type ExecuteResponse struct {
	Body []*pkg.SandboxPhaseResults
}

func Execute(ctx context.Context, input *ExecuteRequest) (*ExecuteResponse, error) {
	sandbox, err := internal.SandboxPool.Acquire()
	if err != nil {
		return nil, err
	}

	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}

	for _, prep := range input.Body.Prepare {
		if err := sandbox.Prepare(&prep); err != nil {
			return nil, err
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

	if err := internal.SandboxPool.Release(sandbox); err != nil {
		return nil, err
	}

	return &ExecuteResponse{Body: results}, nil
}

type ExecuteWithEngineRequest struct {
	Engine string `path:"engine" example:"python"`
	Body   struct {
		Options pkg.SandboxPhaseOptions `json:"options" doc:"Options and limits for the sandbox"`
		Files   []pkg.SandboxFile       `json:"files" doc:"Files to mount in the sandbox before execution"`
	}
}

type ExecuteWithEngineResponse struct {
	Body *pkg.SandboxPhaseResults
}

func ExecuteWithEngine(
	ctx context.Context,
	input *ExecuteWithEngineRequest,
) (*ExecuteWithEngineResponse, error) {
	engine, err := internal.EngineManager.Get(input.Engine)
	if err != nil {
		return nil, err
	}

	sandbox, err := internal.SandboxPool.Acquire()
	if err != nil {
		return nil, err
	}

	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}

	output, err := engine.Run(sandbox, &input.Body.Options)
	if err != nil {
		return nil, err
	}

	if err := internal.SandboxPool.Release(sandbox); err != nil {
		return nil, err
	}

	return &ExecuteWithEngineResponse{Body: output}, nil
}
