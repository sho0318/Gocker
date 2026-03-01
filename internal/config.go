package internal

type Config struct {
	CgroupPath      string
	MaxProcesses    string
	MaxMemory       string
	MaxSwap         string

	RootfsPath      string
	Hostname        string

	BridgeName      string
	BridgeIP        string
	ContainerIP     string
	NetworkBase     string

	Gateway         string
	
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
		BridgeName:   "br0",
		BridgeIP:     "172.18.0.1/24",
		ContainerIP:  "172.18.0.2/24",
		NetworkBase:  "172.18.0",
		Gateway:      "172.18.0.1",
		Command:      command,
	}
}
