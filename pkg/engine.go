package pkg

import (
	"encoding/json"
	"os"
	"path"
)

type Engine struct {
	Compile *SandboxPhase `json:"compile" required:"false"`
	Execute *SandboxPhase `json:"execute"`
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

func (m *EngineManager) List() ([]string, error) {
	entries, err := os.ReadDir(m.Index)
	if err != nil {
		return nil, err
	}

	engines := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if path.Ext(name) != ".json" {
			continue
		}
		engines = append(engines, name[:len(name)-len(".json")])
	}

	return engines, nil
}
