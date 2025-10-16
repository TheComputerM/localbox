package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thecomputerm/localbox/internal"
	"github.com/thecomputerm/localbox/internal/routes"
	"github.com/thecomputerm/localbox/pkg"
)

func init() {
	if os.Getuid() != 0 {
		log.Println("WARNING: LocalBox is not running as root")
	}
	if err := internal.InitCGroup(); err != nil {
		log.Fatal(errors.Join(fmt.Errorf("couldn't init cgroup"), err))
	}
}

func main() {
	huma.DefaultArrayNullable = false
	cli := humacli.New(func(h humacli.Hooks, o *pkg.LocalboxConfig) {
		if err := pkg.SetupLocalbox(o); err != nil {
			log.Fatal(err)
		}

		router := chi.NewMux()
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		config := huma.DefaultConfig("LocalBox", "0.0.7")
		config.Info.Description = `LocalBox is a **easy-to-host**, **general purpose** 
		and **fast** code execution system for running **untrusted** code in sandboxes.
		`
		app := humachi.New(router, config)

		huma.Register(app, huma.Operation{
			OperationID: "system-info",
			Method:      http.MethodGet,
			Path:        "/",
			Summary:     "System Info",
			Description: `Get system information and configuration data.`,
		}, func(ctx context.Context, _ *struct{}) (*routes.SystemInfoResponse, error) {
			return routes.GetSystemInfo(ctx, o)
		})

		huma.Register(app, huma.Operation{
			OperationID: "run-system-gc",
			Method:      http.MethodPost,
			Path:        "/system/gc",
			Summary:     "Run System GC",
			Description: `Run the Nix garbage collector to free up space from unused packages. Will also remove installed engines.`,
		}, routes.RunSystemGC)

		huma.Register(app, huma.Operation{
			OperationID: "execute",
			Method:      http.MethodPost,
			Path:        "/execute",
			Summary:     "Execute Workflow",
			Description: `Execute a series of phases, where each of them can have different options, packages and commands with persistent files. Use this for more complicated workflows.`,
		}, routes.Execute)

		huma.Register(app, huma.Operation{
			OperationID: "list-engines",
			Method:      http.MethodGet,
			Summary:     "List all Engines",
			Path:        "/engine",
			Description: `List all the available engines.`,
			Tags:        []string{"Engine"},
		}, routes.ListEngines)

		huma.Register(app, huma.Operation{
			OperationID: "engine-info",
			Method:      http.MethodGet,
			Path:        "/engine/{engine}",
			Summary:     "Engine Info",
			Description: `Get information about the specified engine.`,
			Tags:        []string{"Engine"},
		}, routes.EngineInfo)

		huma.Register(app, huma.Operation{
			OperationID: "install-engine",
			Method:      http.MethodPost,
			Path:        "/engine/{engine}",
			Summary:     "Install Engine",
			Description: `Install the specified engine.`,
			Tags:        []string{"Engine"},
		}, routes.InstallEngine)

		huma.Register(app, huma.Operation{
			OperationID: "uninstall-engine",
			Method:      http.MethodDelete,
			Path:        "/engine/{engine}",
			Summary:     "Uninstall Engine",
			Description: `Uninstall the specified engine.`,
			Tags:        []string{"Engine"},
		}, routes.UninstallEngine)

		huma.Register(app, huma.Operation{
			OperationID: "execute-engine",
			Method:      http.MethodPost,
			Path:        "/engine/{engine}/execute",
			Summary:     "Execute Engine",
			Description: `Execute a predefined engine with an execution phase whose options can be overriden.`,
			Tags:        []string{"Engine"},
		}, routes.ExecuteEngine)

		h.OnStart(func() {
			fmt.Printf("LocalBox is up and running on :%d\n", o.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", o.Port), router)
		})
	})

	cli.Run()
}
