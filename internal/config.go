package internal

type Config struct {
	CgroupPath      string
	MaxProcesses    string
	MaxMemory       string
	MaxSwap         string

	RootfsPath      string
	Hostname        string
	
	Command         []string
}

func NewDefaultConfig(command []string) *Config {
	return &Config{
		CgroupPath:   "/sys/fs/cgroup/gocker/",
		MaxProcesses: "20",
		MaxMemory:    "10M",
		MaxSwap:      "0",
		RootfsPath:   "./rootfs",
		Hostname:     "gocker-container",
		Command:      command,
	}
}
