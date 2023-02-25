package routes

import (
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
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

func UdpTest(pc net.PacketConn, target string) bool {
	defer func(pc net.PacketConn) {
		err := pc.Close()
		if err != nil {
			log.Warnln("fail to close pc, err: %v\n", err)
		}
	}(pc)
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		log.Warnln("fail to resolve udp, err: %v\n", err)
		return false
	}
	msg, err := MakeDnsMsg(DnsMsg)
	if err != nil {
		log.Warnln("fail to make msg, err: %v\n", err)
		return false
	}
	_, err = pc.WriteTo(msg, addr)
	if err != nil {
		log.Warnln("fail to write to target, err: %v\n", err)
		return false
	}
	err = pc.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Warnln("fail to set deadline for pc, err: %v\n", err)
		return false
	}
	for {
		buf := make([]byte, 1024)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			log.Warnln("fail to read data, err: %v\n", err)
			return false
		}
		log.Infoln("length of dns: %v", n)
		return CheckDnsMsg(buf[:n])
	}
}
