package network

import (
	"strconv"

	"gocker/internal"
)

type Veth struct {
	hostIF string
	containerIF string
	containerPid int
}

func NewVeth(hostIF string, containerIF string, containerPid int) *Veth {
	return &Veth{
		hostIF: hostIF,
		containerIF: containerIF,
		containerPid: containerPid,
	}
}

func (v *Veth) Connect(bridge *Bridge) {
	internal.Run("ip", "link", "add", v.hostIF, "type", "veth", "peer", "name", v.containerIF)
	internal.Run("ip", "link", "set", v.hostIF, "master", bridge.name)
	internal.Run("ip", "link", "set", v.hostIF, "up")
	internal.Run("ip", "link", "set", v.containerIF, "netns", strconv.Itoa(v.containerPid))
}