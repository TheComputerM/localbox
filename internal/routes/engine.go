package routes

import (
	"context"

	"github.com/thecomputerm/localbox/pkg"
)

type ListEnginesResponse struct {
	Body []string `example:"[\"python\", \"go\",\"node\", \"...rest of the engines\"]"`
}

func ListEngines(ctx context.Context, _ *struct{}) (*ListEnginesResponse, error) {
	engines, err := pkg.Globals.EngineManager.List()
	if err != nil {
		return nil, err
	}

	return &ListEnginesResponse{
		Body: engines,
	}, nil
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
	engine, err := pkg.Globals.EngineManager.Get(input.Engine)
	if err != nil {
		return nil, err
	}

	sandbox, err := pkg.Globals.SandboxPool.Acquire()
	if err != nil {
		return nil, err
	}
	defer pkg.Globals.SandboxPool.Release(sandbox)

	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}

	output, err := engine.Run(sandbox, &input.Body.Options)
	if err != nil {
		return nil, err
	}

	return &ExecuteWithEngineResponse{Body: output}, nil
}
