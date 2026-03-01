package host

import (
	"fmt"
	"os"
	"path/filepath"

	"gocker/internal"
)

type Cgroup struct {
	path string
	cfg  *internal.Config
}

func NewCgroup(cfg *internal.Config) *Cgroup {
	return &Cgroup{
		path: cfg.CgroupPath,
		cfg:  cfg,
	}
}

func (cg *Cgroup) Create() {
	internal.Must(os.MkdirAll(cg.path, 0755))
	cg.setResourceLimits()
}

func (cg *Cgroup) AddProcess(pid int) {
	procsPath := filepath.Join(cg.path, "cgroup.procs")
	internal.Must(os.WriteFile(procsPath, []byte(fmt.Sprintf("%d", pid)), 0644))
}

func (cg *Cgroup) Remove() {
	os.Remove(cg.path)
}

func (cg *Cgroup) setResourceLimits() {
	limits := map[string]string{
		"pids.max":        cg.cfg.MaxProcesses,
		"memory.max":      cg.cfg.MaxMemory,
		"memory.swap.max": cg.cfg.MaxSwap,
	}

	for filename, value := range limits {
		path := filepath.Join(cg.path, filename)
		internal.Must(os.WriteFile(path, []byte(value), 0644))
	}
}
