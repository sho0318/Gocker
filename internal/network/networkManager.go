package network

import (
	"fmt"
	"strconv"
	"os"
    "path/filepath"
	"gocker/internal"
	"strings"
)

const ipDbPath = "/var/run/gocker/ips"

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

func (nm *NetworkManager) AllocateIP() string {
	if err := os.MkdirAll(ipDbPath, 0755); err != nil {
        panic(fmt.Errorf("failed to create IP database directory: %v", err))
    }

	for i := 2; i < 254; i++ {
		ip := fmt.Sprintf("172.18.0.%d", i)
		recordFile := filepath.Join(ipDbPath, ip)

		if _, err := os.Stat(recordFile); os.IsNotExist(err) {
            if err := os.WriteFile(recordFile, []byte(""), 0644); err != nil {
                panic(err)
            }
            return ip
        }
	}
	panic("No more IP address available")
}

func (nm *NetworkManager) ReleaseIP(ip string) {
    cleanIP := strings.Split(ip, "/")[0]
    recordFile := filepath.Join(ipDbPath, cleanIP)
    
    if err := os.Remove(recordFile); err != nil {
        fmt.Printf("Failed to release IP: %v\n", err)
    }
}