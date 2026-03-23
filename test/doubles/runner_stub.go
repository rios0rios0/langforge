package doubles

import "os/exec"

// RunnerStub is a test double for cmdexec.Runner.
type RunnerStub struct {
	RunFunc       func(dir, name string, args ...string) error
	RunOutputFunc func(dir, name string, args ...string) (string, error)
}

// NewRunnerStubBinaryNotFound creates a RunnerStub that simulates a missing binary.
func NewRunnerStubBinaryNotFound() *RunnerStub {
	return &RunnerStub{
		RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
			return "", exec.ErrNotFound
		},
	}
}

func (r *RunnerStub) Run(dir string, name string, args ...string) error {
	if r.RunFunc != nil {
		return r.RunFunc(dir, name, args...)
	}
	return nil
}

func (r *RunnerStub) RunOutput(dir string, name string, args ...string) (string, error) {
	if r.RunOutputFunc != nil {
		return r.RunOutputFunc(dir, name, args...)
	}
	return "", nil
}
