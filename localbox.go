package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/labstack/echo/v4"
	"github.com/thecomputerm/localbox/internal/routes"
	"github.com/thecomputerm/localbox/pkg"
)

func main() {
	huma.DefaultArrayNullable = false
	cli := humacli.New(func(h humacli.Hooks, o *pkg.LocalboxConfig) {
		if err := pkg.SetupLocalbox(o); err != nil {
			log.Fatal(err)
		}

		e := echo.New()
		config := huma.DefaultConfig("LocalBox", "1.0.0")
		config.Info.Description = `LocalBox is a **easy-to-host**, **general purpose** 
		and **fast** code execution system for running **untrusted** code in sandboxes.
		`
		app := humaecho.New(e, config)

		huma.Register(app, huma.Operation{
			OperationID: "system-info",
			Method:      http.MethodGet,
			Path:        "/system",
			Summary:     "System Info",
			Description: `Get system information and configuration data.`,
		}, func(ctx context.Context, _ *struct{}) (*routes.SystemInfoResponse, error) {
			return routes.GetSystemInfo(ctx, o)
		})

		huma.Register(app, huma.Operation{
			OperationID: "list-engines",
			Method:      http.MethodGet,
			Summary:     "List Engines",
			Path:        "/engines",
			Description: `List all the available engines.`,
		}, routes.ListEngines)

		huma.Register(app, huma.Operation{
			OperationID: "run-engine",
			Method:      http.MethodPost,
			Path:        "/engine/{engine}",
			Summary:     "Run Engine",
			Description: `Execute a predefined engine with an execution phase whose options can be overriden.`,
		}, routes.ExecuteWithEngine)

		huma.Register(app, huma.Operation{
			OperationID: "execute",
			Method:      http.MethodPost,
			Path:        "/execute",
			Summary:     "Execute",
			Description: `Execute a series of phases, where each of them can have different options, packages and commands with persistent files. Use this for more complicated workflows.`,
		}, routes.Execute)

		h.OnStart(func() {
			e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", o.Port)))
		})
	})

	cli.Run()
}
