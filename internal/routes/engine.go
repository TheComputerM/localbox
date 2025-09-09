package routes

import (
	"context"

	"github.com/thecomputerm/localbox/internal"
)

type ListEnginesResponse struct {
	Body []string `example:"[\"python\", \"go\",\"node\", \"...rest of the engines\"]"`
}

func ListEngines(ctx context.Context, _ *struct{}) (*ListEnginesResponse, error) {
	engines, err := internal.EngineManager.List()
	if err != nil {
		return nil, err
	}

	return &ListEnginesResponse{
		Body: engines,
	}, nil
}
