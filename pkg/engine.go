package pkg

import (
	"encoding/json"
	"os"
	"path"
)

type EngineMetadata struct {
	Version string `json:"version"`
}

type Engine struct {
	Compile *SandboxPrepare `json:"compile,omitempty" required:"false"`
	Execute *SandboxPhase   `json:"execute"`
	Meta    *EngineMetadata `json:"meta"`
}

type EngineManager struct {
	// Path to predefined engine definitions
	Index string
}

func (m *EngineManager) Get(name string) (*Engine, error) {
	content, err := os.ReadFile(path.Join(m.Index, name+".json"))
	if err != nil {
		return nil, err
	}

	var engine Engine
	err = json.Unmarshal(content, &engine)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func (m *EngineManager) List() (map[string]EngineMetadata, error) {
	entries, err := os.ReadDir(m.Index)
	if err != nil {
		return nil, err
	}

	engines := make(map[string]EngineMetadata, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if path.Ext(name) != ".json" {
			continue
		}
		name = name[:len(name)-len(".json")]
		engine, err := m.Get(name)
		if err != nil {
			return nil, err
		}
		engines[name] = *engine.Meta
	}

	return engines, nil
}

func (e *Engine) Run(s Sandbox, options *SandboxPhaseOptions) (*SandboxPhaseResults, error) {
	if e.Compile != nil {
		if err := s.Prepare(e.Compile); err != nil {
			return nil, err
		}
	}

	result, err := s.Run(e.Execute, options)
	if err != nil {
		return nil, err
	}

	return result, nil
}
