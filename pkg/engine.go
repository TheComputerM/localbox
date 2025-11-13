package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
)

type EngineMetadata struct {
	RunFile string `json:"run_file" doc:"Name of used to replace the '@' file when the engine is executed" example:"main.py"`
	Version string `json:"version" doc:"Version of the runtime or compiler used by the engine" example:"3.12.11" readOnly:"true"`
}

type EngineInfo struct {
	EngineMetadata
	Installed bool `json:"installed" doc:"Whether the packages used by the engine are installed" example:"false"`
}

type Engine struct {
	Compile *SandboxCommand `json:"compile" doc:"Optional compilation phase before execution" required:"false"`
	Execute *SandboxPhase   `json:"execute" doc:"Main execution phase"`
	Meta    *EngineMetadata `json:"meta" doc:"Metadata about the engine"`
}

func (e *Engine) packages() []string {
	packages := e.Execute.Packages
	if e.Compile != nil {
		packages = append(packages, e.Compile.Packages...)
	}
	slices.Sort(packages)
	return slices.Compact(packages)
}

func (e *Engine) isInstalled() bool {
	args := []string{"path-info"}
	args = append(args, e.packages()...)
	cmd := exec.Command("nix", args...)
	return cmd.Run() == nil
}

func (e *Engine) Install() error {
	args := []string{"shell"}
	args = append(args, e.packages()...)
	args = append(args, "-c", "true")
	cmd := exec.Command("nix", args...)
	return cmd.Run()
}

func (e *Engine) Uninstall() error {
	args := []string{"store", "delete"}
	args = append(args, e.packages()...)
	cmd := exec.Command("nix", args...)
	return cmd.Run()
}

func (e *Engine) Info() *EngineInfo {
	return &EngineInfo{
		EngineMetadata: *e.Meta,
		Installed:      e.isInstalled(),
	}
}

func (e *Engine) Run(s Sandbox, options *SandboxPhaseOptions) (*SandboxPhaseResults, error) {
	if e.Compile != nil {
		if result, ok := s.UnsafeRun(e.Compile); !ok {
			return result, nil
		}
	}

	result, err := s.Run(e.Execute, options)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type EngineManager struct {
	// Path to predefined engine definitions
	Root string
}

func (m *EngineManager) Get(name string) (*Engine, error) {
	content, err := os.ReadFile(path.Join(m.Root, name+".json"))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("engine %s not found", name), err)
	}

	var engine Engine
	err = json.Unmarshal(content, &engine)
	if err != nil {
		return nil, err
	}

	if engine.Execute == nil {
		return nil, fmt.Errorf("%s engine doesn't have a execute field", name)
	}

	return &engine, nil
}

func (m *EngineManager) List() (map[string]*EngineInfo, error) {
	entries, err := os.ReadDir(m.Root)
	if err != nil {
		return nil, err
	}

	engines := make(map[string]*EngineInfo, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if path.Ext(name) != ".json" {
			continue
		}
		name = strings.TrimSuffix(name, ".json")
		engine, err := m.Get(name)
		if err != nil {
			return nil, err
		}
		engines[name] = engine.Info()
	}

	return engines, nil
}
