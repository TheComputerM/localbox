package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/thecomputerm/localbox/pkg"
)

type ListEnginesResponse struct {
	Body map[string]*pkg.EngineInfo `example:"{\"python\": {\"version\": \"3.12.11\", \"run_file\": \"main.py\", \"installed\": true}}"`
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

type EngineRequest struct {
	Engine string `path:"engine" example:"python"`
}

type EngineInfoResponse struct {
	Body *pkg.EngineInfo `example:"{\"version\": \"3.12.11\", \"run_file\": \"main.py\", \"installed\": true}"`
}

func EngineInfo(ctx context.Context, input *EngineRequest) (*EngineInfoResponse, error) {
	engine, err := pkg.Globals.EngineManager.Get(input.Engine)
	if err != nil {
		return nil, err
	}
	return &EngineInfoResponse{Body: engine.Info()}, nil
}

func InstallEngine(ctx context.Context, input *EngineRequest) (*struct{}, error) {
	engine, err := pkg.Globals.EngineManager.Get(input.Engine)
	if err != nil {
		return nil, err
	}
	if err := engine.Install(); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to install %s engine", input.Engine), err)
	}
	return nil, nil
}

func UninstallEngine(ctx context.Context, input *EngineRequest) (*struct{}, error) {
	engine, err := pkg.Globals.EngineManager.Get(input.Engine)
	if err != nil {
		return nil, err
	}
	if err := engine.Uninstall(); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to uninstall %s engine", input.Engine), err)
	}
	return nil, nil
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

	for i := range input.Body.Files {
		if input.Body.Files[i].Name == "@" {
			input.Body.Files[i].Name = engine.Meta.RunFile
		}
	}

	if err := sandbox.Mount(input.Body.Files); err != nil {
		return nil, err
	}

	output, err := engine.Run(sandbox, &input.Body.Options)
	if err != nil {
		return nil, err
	}

	return &ExecuteWithEngineResponse{Body: output}, nil
}
