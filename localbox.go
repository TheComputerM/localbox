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

		routes.AddRoutes(app)

		h.OnStart(func() {
			fmt.Printf("LocalBox is up and running on :%d\n", o.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", o.Port), router)
		})
	})

	cli.Run()
}
