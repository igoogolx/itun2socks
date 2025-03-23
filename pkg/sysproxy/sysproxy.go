package sysproxy

import "net"

func Set(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	err = SetWebProxy(host, port)
	if err != nil {
		return err
	}
	err = SetSOCKSProxy(host, port)
	if err != nil {
		return err
	}
	return nil
}

func Clear() error {
	err := DisableWebProxy()
	if err != nil {
		return err
	}
	err = DisableSOCKSProxy()
	if err != nil {
		return err
	}
	return nil
}
