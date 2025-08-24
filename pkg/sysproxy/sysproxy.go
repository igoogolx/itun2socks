package sysproxy

import "net"

func Set(addr string, activeInterface string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	return SetWebProxy(host, port, activeInterface)
}

func Clear(activeInterface string) error {
	return DisableWebProxy(activeInterface)
}
