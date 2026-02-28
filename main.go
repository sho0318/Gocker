package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	//Cgroup作成
	cgroupPath := "/sys/fs/cgroup/gocker/"
	setControlGroup(cgroupPath)

	//子プロセス準備
	cmd := exec.Command("./gocker", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	//子プロセス実行
	must(cmd.Start())
	pid := strconv.Itoa(cmd.Process.Pid)
	must(os.WriteFile(filepath.Join(cgroupPath, "cgroup.procs"), []byte(pid), 0644))
	must(cmd.Wait())

	//Cgroup削除
	removeControlGroup(cgroupPath)
}

func setControlGroup(cgroupPath string) {
	must(os.MkdirAll(cgroupPath, 0755))
	must(os.WriteFile(filepath.Join(cgroupPath, "pids.max"), []byte("20"), 0700))
	must(os.WriteFile(filepath.Join(cgroupPath, "memory.max"), []byte("10M"), 0644))
	must(os.WriteFile(filepath.Join(cgroupPath, "memory.swap.max"), []byte("0"), 0644))
}

func removeControlGroup(cgroupPath string) {
	must(os.Remove(cgroupPath))
}

func child() {
	fmt.Printf("Running %v \n", os.Args[2:])

	must(syscall.Chroot("./rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(syscall.Sethostname([]byte("gocker-container")))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}