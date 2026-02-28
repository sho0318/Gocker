package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

type Cgroup struct {
	path string
	cfg  *Config
}

func NewCgroup(cfg *Config) *Cgroup {
	return &Cgroup{
		path: cfg.CgroupPath,
		cfg:  cfg,
	}
}

func (cg *Cgroup) Create() error {
	if err := os.MkdirAll(cg.path, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory: %w", err)
	}

	if err := cg.setResourceLimits(); err != nil {
		return fmt.Errorf("failed to set resource limits: %w", err)
	}

	return nil
}

func (cg *Cgroup) AddProcess(pid int) error {
	procsPath := filepath.Join(cg.path, "cgroup.procs")
	if err := os.WriteFile(procsPath, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %w", err)
	}
	return nil
}

func (cg *Cgroup) Remove() error {
	if err := os.Remove(cg.path); err != nil {
		return fmt.Errorf("failed to remove cgroup: %w", err)
	}
	return nil
}

func (cg *Cgroup) setResourceLimits() error {
	limits := map[string]string{
		"pids.max":        cg.cfg.MaxProcesses,
		"memory.max":      cg.cfg.MaxMemory,
		"memory.swap.max": cg.cfg.MaxSwap,
	}

	for filename, value := range limits {
		path := filepath.Join(cg.path, filename)
		if err := os.WriteFile(path, []byte(value), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}
