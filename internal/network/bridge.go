package network

import (
	"os/exec"
	"strings"

	"gocker/internal"
)

type Bridge struct {
	name string
	ip string
}

func NewBridge(config *internal.Config) *Bridge {
	bridgeName := config.BridgeName
	bridgeIP := config.BridgeIP
	return &Bridge{
		name: bridgeName,
		ip: bridgeIP,
	}
}

func (b *Bridge) Create() {
	ensureBridgeExists(b.name, b.ip)
	internal.Run("sysctl", "-w", "net.ipv4.ip_forward=1")
	setupNAT(b.name, b.ip)
}

func ensureBridgeExists(
	bridgeName string,
	bridgeIP string,
) {
	output, err := exec.Command("ip", "link", "show", bridgeName).CombinedOutput()
	bridgeExists := err == nil && strings.Contains(string(output), bridgeName)

	if !bridgeExists {
		internal.Run("ip", "link", "add", "name", bridgeName, "type", "bridge")
	}

	output = internal.RunOutput("ip", "link", "show", bridgeName)

	if !strings.Contains(string(output), "UP") {
		internal.Run("ip", "link", "set", bridgeName, "up")
	}

	addrOutput, _ := exec.Command("ip", "addr", "show", bridgeName).CombinedOutput()
	if !strings.Contains(string(addrOutput), strings.Split(bridgeIP, "/")[0]) {
		internal.Run("ip", "addr", "add", bridgeIP, "dev", bridgeName)
	}
}

func setupNAT(bridgeName string, bridgeIP string) {
    subnet := strings.Split(bridgeIP, ".")[0] + "." + 
              strings.Split(bridgeIP, ".")[1] + "." + 
              strings.Split(bridgeIP, ".")[2] + ".0/24"

    exec.Command("iptables", "-t", "nat", "-D", "POSTROUTING", "-s", subnet, "!", "-o", bridgeName, "-j", "MASQUERADE").Run()

    internal.Run("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", subnet, "!", "-o", bridgeName, "-j", "MASQUERADE")

    internal.Run("iptables", "-A", "FORWARD", "-i", bridgeName, "-j", "ACCEPT")
    internal.Run("iptables", "-A", "FORWARD", "-o", bridgeName, "-j", "ACCEPT")
}