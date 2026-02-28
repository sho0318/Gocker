package internal

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Container struct {
	cfg *Config
	cg  *Cgroup
}

func NewContainer(cfg *Config) *Container {
	return &Container{
		cfg: cfg,
		cg:  NewCgroup(cfg),
	}
}

func (c *Container) Run() error {
	if err := c.cg.Create(); err != nil {
		return fmt.Errorf("failed to create cgroup: %w", err)
	}
	defer func() {
		if err := c.cg.Remove(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove cgroup: %v\n", err)
		}
	}()

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, c.cfg.Command...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start child process: %w", err)
	}

	if err := c.cg.AddProcess(cmd.Process.Pid); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %w", err)
	}

	cmd.Wait()
	return nil
}

func (c *Container) RunChild() error {
	fmt.Printf("Running %v\n", c.cfg.Command)

	if err := syscall.Chroot(c.cfg.RootfsPath); err != nil {
		return fmt.Errorf("failed to chroot: %w", err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to chdir to /: %w", err)
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("failed to mount proc: %w", err)
	}
	defer func() {
		if err := syscall.Unmount("proc", 0); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to unmount proc: %v\n", err)
		}
	}()

	if err := syscall.Sethostname([]byte(c.cfg.Hostname)); err != nil {
		return fmt.Errorf("failed to set hostname: %w", err)
	}

	cmd := exec.Command(c.cfg.Command[0], c.cfg.Command[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	return cmd.Run()
}
