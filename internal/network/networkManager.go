package network

import (
	"fmt"
	"strconv"

	"gocker/internal"
)

type NetworkManager struct {
	config *internal.Config
	bridge *Bridge
}

func NewNetworkManager(config *internal.Config) *NetworkManager {
	return &NetworkManager{
		config: config,
		bridge: NewBridge(config),
	}
}

func (nm *NetworkManager) Setup() {
	nm.bridge.Create()
}

func (nm *NetworkManager) Connect(
	containerPid int,
) {
	hostIFName := fmt.Sprintf("veth%d", containerPid)
	containerIFName := fmt.Sprintf("vethc%d", containerPid)
	veth := NewVeth(hostIFName, containerIFName, containerPid)

	veth.Connect(nm.bridge)

	nm.ConfigureNamespace(containerPid, containerIFName)
}

func (nm *NetworkManager) ConfigureNamespace(pid int, ifName string) {
    pidStr := strconv.Itoa(pid)

    internal.Run("nsenter", "-t", pidStr, "-n", "ip", "link", "set", "lo", "up")
    internal.Run("nsenter", "-t", pidStr, "-n", "ip", "link", "set", ifName, "name", "eth0")
    internal.Run("nsenter", "-t", pidStr, "-n", "ip", "addr", "add", nm.config.ContainerIP, "dev", "eth0")
    internal.Run("nsenter", "-t", pidStr, "-n", "ip", "link", "set", "eth0", "up")
    internal.Run("nsenter", "-t", pidStr, "-n", "ip", "route", "add", "default", "via", nm.config.Gateway, "dev", "eth0")

    internal.Run("nsenter", "-t", pidStr, "-n", "sysctl", "-w", "net.ipv4.conf.all.accept_redirects=0")
    internal.Run("nsenter", "-t", pidStr, "-n", "sysctl", "-w", "net.ipv4.conf.eth0.accept_redirects=0")
}