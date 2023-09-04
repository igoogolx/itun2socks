package routes

import (
	"github.com/igoogolx/itun2socks/pkg/log"
	"github.com/miekg/dns"
	"net"
	"time"
)

const DnsMsg = "google.com."

func MakeDnsMsg(host string) ([]byte, error) {
	m := new(dns.Msg)
	m.SetQuestion(host, dns.TypeA)
	return m.Pack()
}

func CheckDnsMsg(data []byte) bool {
	m := new(dns.Msg)
	err := m.Unpack(data)
	if err != nil {
		return false
	}
	return true
}

func UdpTest(pc net.PacketConn, target string) (bool, error) {
	defer func(pc net.PacketConn) {
		err := pc.Close()
		if err != nil {
			log.Debugln(log.FormatLog(log.HubPrefix, "fail to close conn in udp test, err: %v"), err)
		}
	}(pc)
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return false, err
	}
	msg, err := MakeDnsMsg(DnsMsg)
	if err != nil {
		return false, err
	}
	_, err = pc.WriteTo(msg, addr)
	if err != nil {
		return false, err
	}
	err = pc.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return false, err
	}
	for {
		buf := make([]byte, 1024)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			return false, err
		}
		return CheckDnsMsg(buf[:n]), nil
	}
}
