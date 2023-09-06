package tun

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	F "github.com/sagernet/sing/common/format"
	"net"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	Name      string
	LocalAddr *net.IPNet
	Gateway   *net.IPNet
	Mtu       int
}

func CalculateInterfaceName(name string) (tunName string) {
	if runtime.GOOS == "darwin" {
		tunName = "utun"
	} else if name != "" {
		tunName = name
		return
	} else {
		tunName = "tun"
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	var tunIndex int
	for _, netInterface := range interfaces {
		if strings.HasPrefix(netInterface.Name, tunName) {
			index, parseErr := strconv.ParseInt(netInterface.Name[len(tunName):], 10, 16)
			if parseErr == nil {
				tunIndex = int(index) + 10
			}
		}
	}
	tunName = F.ToString(tunName, tunIndex)
	return
}

func New() (Config, error) {
	local, gw, err := ParseAddresses(
		constants.TunLocalAddr,
		constants.TunMask,
		constants.TunGateway,
	)
	if err != nil {
		return Config{}, err
	}
	tunInterfaceName := CalculateInterfaceName(constants.TunName)
	return Config{
		Name:      tunInterfaceName,
		Mtu:       constants.TunMtu,
		LocalAddr: local,
		Gateway:   gw,
	}, nil
}
