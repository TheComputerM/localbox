package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/labstack/echo/v4"
	"github.com/thecomputerm/localbox/internal/routes"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"9000"`
}

func main() {
	if os.Getuid() != 0 {
		log.Fatal("LocalBox must be run as root")
	}

	huma.DefaultArrayNullable = false
	cli := humacli.New(func(h humacli.Hooks, o *Options) {
		e := echo.New()
		config := huma.DefaultConfig("LocalBox", "1.0.0")
		config.Info.Description = `LocalBox is a **easy-to-host**, **general purpose** 
		and **fast** code execution system for running **untrusted** code in sandboxes.
		`
		app := humaecho.New(e, config)

		huma.Register(app, huma.Operation{
			OperationID: "list-engines",
			Method:      http.MethodGet,
			Summary:     "List Engines",
			Path:        "/engines",
			Description: `List all the available engines.`,
		}, routes.ListEngines)

		huma.Register(app, huma.Operation{
			OperationID: "execute",
			Method:      http.MethodPost,
			Path:        "/execute",
			Summary:     "Execute",
			Description: `Execute a series of phases, where each of them can have different options, packages and commands with persistent files.`,
		}, routes.Execute)

		huma.Register(app, huma.Operation{
			OperationID: "execute-engine",
			Method:      http.MethodPost,
			Path:        "/execute/{engine}",
			Summary:     "Execute Engine",
			Description: `Execute a predefined engine that has execute and (an optional) compile phase whose options can be overriden.`,
		}, routes.ExecuteWithEngine)

		h.OnStart(func() {
			e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", o.Port)))
		})
	})

	cli.Run()
}
