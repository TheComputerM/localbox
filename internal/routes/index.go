package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func AddRoutes(app huma.API) {
	huma.Register(app, huma.Operation{
		OperationID: "run-system-gc",
		Method:      http.MethodPost,
		Path:        "/system/gc",
		Summary:     "Run System GC",
		Description: `Run the Nix garbage collector to free up space from unused packages. Will also remove installed engines.`,
	}, RunSystemGC)

	huma.Register(app, huma.Operation{
		OperationID: "execute",
		Method:      http.MethodPost,
		Path:        "/execute",
		Summary:     "Execute Workflow",
		Description: `Execute a series of phases, where each of them can have different options, packages and commands with persistent files. Use this for more complicated workflows.`,
	}, Execute)

	huma.Register(app, huma.Operation{
		OperationID: "list-engines",
		Method:      http.MethodGet,
		Summary:     "List all Engines",
		Path:        "/engine",
		Description: `List all the available engines.`,
		Tags:        []string{"Engine"},
	}, ListEngines)

	huma.Register(app, huma.Operation{
		OperationID: "engine-info",
		Method:      http.MethodGet,
		Path:        "/engine/{engine}",
		Summary:     "Engine Info",
		Description: `Get information about the specified engine.`,
		Tags:        []string{"Engine"},
	}, EngineInfo)

	huma.Register(app, huma.Operation{
		OperationID: "install-engine",
		Method:      http.MethodPost,
		Path:        "/engine/{engine}",
		Summary:     "Install Engine",
		Description: `Install the specified engine.`,
		Tags:        []string{"Engine"},
	}, InstallEngine)

	huma.Register(app, huma.Operation{
		OperationID: "uninstall-engine",
		Method:      http.MethodDelete,
		Path:        "/engine/{engine}",
		Summary:     "Uninstall Engine",
		Description: `Uninstall the specified engine.`,
		Tags:        []string{"Engine"},
	}, UninstallEngine)

	huma.Register(app, huma.Operation{
		OperationID: "execute-engine",
		Method:      http.MethodPost,
		Path:        "/engine/{engine}/execute",
		Summary:     "Execute Engine",
		Description: `Execute a predefined engine with an execution phase whose options can be overriden.`,
		Tags:        []string{"Engine"},
	}, ExecuteEngine)

	huma.Register(app, huma.Operation{
		OperationID: "execute-custom-engine",
		Method:      http.MethodPost,
		Path:        "/engine/custom",
		Summary:     "Execute Custom Engine",
		Description: `Execute a custom engine with an execution phase whose options can be overridden.`,
		Tags:        []string{"Engine"},
	}, ExecuteCustomEngine)
}
