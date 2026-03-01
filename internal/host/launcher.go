package host

import (
	"os"
    "os/exec"
    "syscall"
	"strings"
	"fmt"

	"gocker/internal"
	"gocker/internal/network"
	"gocker/internal/filesystem"
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

	rawIP := nm.AllocateIP()
	defer nm.ReleaseIP(rawIP)
	l.config.ContainerIP = rawIP + "/24"

	ipParts := strings.Split(rawIP, ".")
    id := ipParts[len(ipParts)-1]
	l.config.ContainerID = id

	cmd := exec.Command("/proc/self/exe", append([]string{"child", l.config.ContainerID}, l.config.Command...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	internal.Must(cmd.Start())

	nm.Connect(cmd.Process.Pid)

	l.cgroup.AddProcess(cmd.Process.Pid)

	internal.Must(cmd.Wait())

	om := filesystem.NewOverlayManager(l.config.RootfsPath, l.config.RuntimePath, l.config.ContainerID)
	if err := om.UnmountAndCleanup(); err != nil {
		fmt.Printf("Warning: cleanup failed: %v\n", err)
	}
}