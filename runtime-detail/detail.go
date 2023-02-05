package runtime_detail

import (
	"github.com/Dreamacro/clash/component/iface"
	"github.com/igoogolx/itun2socks/global"
	"sync"
)

type Detail struct {
	DirectedInterfaceName   string   `json:"directedInterfaceName"`
	DirectedInterfaceV4Addr string   `json:"directedInterfaceV4Addr"`
	TunInterfaceName        string   `json:"tunInterfaceName"`
	LocalDns                []string `json:"localDns"`
	RemoteDns               []string `json:"remoteDns"`
}

var defaultRuntimeDetail = Detail{}
var mutex sync.RWMutex

func Update(tunInterfaceName string, localDns []string, remoteDns []string) error {
	mutex.Lock()
	defer mutex.Unlock()
	networkInterface, err := iface.ResolveInterface(global.GetDefaultInterfaceName())
	if err != nil {
		return err
	}
	addr, err := networkInterface.PickIPv4Addr(nil)
	if err != nil {
		return err
	}
	defaultRuntimeDetail = Detail{
		DirectedInterfaceV4Addr: addr.IP.String(),
		DirectedInterfaceName:   networkInterface.Name,
		TunInterfaceName:        tunInterfaceName,
		LocalDns:                localDns,
		RemoteDns:               remoteDns,
	}
	return nil
}

func Get() Detail {
	mutex.RLock()
	defer mutex.RUnlock()
	return defaultRuntimeDetail
}
