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
		Phases []ExecutePhase    `json:"phases"`
		Files  []pkg.SandboxFile `json:"files"`
	}
}

type ExecuteResponse struct {
	Body []*pkg.SandboxPhaseResults
}

func Execute(ctx context.Context, input *ExecuteRequest) (*ExecuteResponse, error) {
	sandbox := internal.SandboxPool.Acquire()

	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}
	results := make([]*pkg.SandboxPhaseResults, len(input.Body.Phases))

	for i, phase := range input.Body.Phases {
		result, err := sandbox.Run(&phase.SandboxPhase, &phase.SandboxPhaseOptions)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	if err := sandbox.Cleanup(); err != nil {
		return nil, err
	}

	return &ExecuteResponse{Body: results}, nil
}

type ExecuteWithEngineRequest struct {
	Engine string `path:"engine"`
	Body   struct {
		Compile pkg.SandboxPhaseOptions `json:"compile,omitempty"`
		Execute pkg.SandboxPhaseOptions `json:"execute"`
		Files   []pkg.SandboxFile       `json:"files"`
	}
}

type ExecuteWithEngineResponse struct {
	Body struct {
		Compile *pkg.SandboxPhaseResults `json:"compile,omitempty"`
		Execute *pkg.SandboxPhaseResults `json:"execute"`
	}
}

func ExecuteWithEngine(
	ctx context.Context,
	input *ExecuteWithEngineRequest,
) (*ExecuteWithEngineResponse, error) {
	engine, err := internal.EngineManager.Get(input.Engine)
	if err != nil {
		return nil, err
	}

	sandbox := internal.SandboxPool.Acquire()
	defer internal.SandboxPool.Release(sandbox)

	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}

	output := &ExecuteWithEngineResponse{}

	if engine.Compile != nil {
		result, err := sandbox.Run(engine.Compile, &input.Body.Compile)
		if err != nil {
			return nil, err
		}
		output.Body.Compile = result
	}

	result, err := sandbox.Run(engine.Execute, &input.Body.Execute)
	if err != nil {
		return nil, err
	}
	output.Body.Execute = result

	return output, nil
}
