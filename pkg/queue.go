package pkg

import (
	"os"
	"os/exec"
)

type SandboxPool struct {
	sandboxes chan Sandbox
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

func (p *SandboxPool) InitCGroup() error {
	if _, err := os.Stat("/sys/fs/cgroup/isolate"); os.IsExist(err) {
		return nil
	}
	cmd := exec.Command("bash", "-c", "mkdir -p isolate && echo '+cpuset +memory' > isolate/cgroup.subtree_control")
	cmd.Dir = "/sys/fs/cgroup"
	return cmd.Run()
}

func (p *SandboxPool) Acquire() Sandbox {
	sandbox := <-p.sandboxes
	sandbox.Init()
	return sandbox
}

func (p *SandboxPool) Release(s Sandbox) error {
	if err := s.Cleanup(); err != nil {
		return err
	}
	p.sandboxes <- s
	return nil
}
