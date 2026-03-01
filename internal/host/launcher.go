package host

import (
	"os"
    "os/exec"
    "syscall"

	"gocker/internal"
	"gocker/internal/network"
)

type Launcher struct {
	config *internal.Config
	cgroup *Cgroup
}

func NewLauncher(config *internal.Config) *Launcher {
	return &Launcher{
		config: config,
		cgroup: NewCgroup(config),
	}
}

func (l *Launcher) Start() {
	l.cgroup.Create()
	defer l.cgroup.Remove()

	nm := network.NewNetworkManager(l.config)
	nm.Setup()

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, l.config.Command...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	internal.Must(cmd.Start())

	nm.Connect(cmd.Process.Pid)

	l.cgroup.AddProcess(cmd.Process.Pid)

	internal.Must(cmd.Wait())
}