package sysproxy

import "net"

func Set(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	return SetWebProxy(host, port)
}

func Clear() error {
	return DisableWebProxy()
}
