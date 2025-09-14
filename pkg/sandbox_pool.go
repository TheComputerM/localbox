package pkg

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

// Gets a new sandbox from the pool and initializes it
func (p *SandboxPool) Acquire() (Sandbox, error) {
	sandbox := <-p.sandboxes
	if err := sandbox.Init(); err != nil {
		return -1, err
	}
	return sandbox, nil
}

// Cleans up the sandbox and returns it to the pool
func (p *SandboxPool) Release(s Sandbox) error {
	if err := s.Cleanup(); err != nil {
		return err
	}
	p.sandboxes <- s
	return nil
}

// Number of available sandboxes
func (p *SandboxPool) Available() int {
	return len(p.sandboxes)
}

func (p *SandboxPool) Capacity() int {
	return cap(p.sandboxes)
}
