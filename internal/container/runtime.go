package container

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"

    "gocker/internal"
)

type Runtime struct {
    config *internal.Config
}

func NewRuntime(config *internal.Config) *Runtime {
    return &Runtime{config: config}
}

func (r *Runtime) Start() error {
	fmt.Printf("Running %v\n", r.config.Command)
    syscall.Sethostname([]byte(r.config.Hostname))

    syscall.Chroot(r.config.RootfsPath)
    os.Chdir("/")

    syscall.Mount("proc", "proc", "proc", 0, "")

	if err := r.setupDNS(); err != nil {
        fmt.Printf("Warning: failed to setup DNS: %v\n", err)
    }

    cmd := exec.Command(r.config.Command[0], r.config.Command[1:]...)
    cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

    return cmd.Run()
}

func (r *Runtime) setupDNS() error {
    if err := os.MkdirAll("/etc", 0755); err != nil {
        return err
    }
    return os.WriteFile("/etc/resolv.conf", []byte("nameserver 8.8.8.8\n"), 0644)
}