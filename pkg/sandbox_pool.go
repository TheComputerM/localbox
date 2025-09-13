package pkg

import (
	"errors"
	"os"
	"path/filepath"
)

type SandboxPool struct {
	sandboxes chan Sandbox
}

func (p *SandboxPool) InitCGroup() error {
	cgroups := "/sys/fs/cgroup"
	name := "isolate"

	if info, err := os.Stat(filepath.Join(cgroups, name)); !errors.Is(err, os.ErrNotExist) && info.IsDir() {
		return nil
	}

	if err := os.Mkdir(filepath.Join(cgroups, name), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(cgroups, name, "cgroup.subtree_control"), []byte("+cpuset +memory"), 0644); err != nil {
		return err
	}
	return nil
}

func NewSandboxPool(size int) *SandboxPool {
	pool := &SandboxPool{
		sandboxes: make(chan Sandbox, size),
	}
	for i := 0; i < size; i++ {
		pool.sandboxes <- Sandbox(i)
	}
	return pool
}

func (p *SandboxPool) Acquire() (Sandbox, error) {
	sandbox := <-p.sandboxes
	if err := sandbox.Init(); err != nil {
		return -1, err
	}
	return sandbox, nil
}

func (p *SandboxPool) Release(s Sandbox) error {
	if err := s.Cleanup(); err != nil {
		return err
	}
	p.sandboxes <- s
	return nil
}
