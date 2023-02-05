package tunnel

import (
	"errors"
	"github.com/Dreamacro/clash/constant"
	"net"
	"strconv"
)

func CreateUdpMetadata(srcAddr, destAddr net.UDPAddr) constant.Metadata {
	metadata := constant.Metadata{
		SrcIP:   srcAddr.IP,
		SrcPort: strconv.Itoa(srcAddr.Port),
		DstIP:   destAddr.IP,
		DstPort: strconv.Itoa(destAddr.Port),
		NetWork: constant.UDP,
	}
	return metadata
}

func CreateTcpMetadata(srcAddr, destAddr net.TCPAddr) constant.Metadata {
	metadata := constant.Metadata{
		SrcIP:   srcAddr.IP,
		SrcPort: strconv.Itoa(srcAddr.Port),
		DstIP:   destAddr.IP,
		DstPort: strconv.Itoa(destAddr.Port),
		NetWork: constant.TCP,
	}
	return metadata
}

func CreateMetadata(srcAddr, destAddr string, network constant.NetWork) (*constant.Metadata, error) {
	var srcHost, srcPort string
	var srcIp net.IP
	var err error
	if len(srcAddr) != 0 {
		srcHost, srcPort, err = net.SplitHostPort(srcAddr)
		if err != nil {
			return nil, err
		}
		srcIp = net.ParseIP(srcHost)
		if srcIp == nil {
			return nil, errors.New("fail to parse src host")
		}
	}

	destHost, destPort, err := net.SplitHostPort(destAddr)

	destIp := net.ParseIP(destHost)
	if destIp == nil {
		return nil, errors.New("fail to parse dest host")
	}
	metadata := &constant.Metadata{
		SrcIP:   srcIp,
		SrcPort: srcPort,
		DstIP:   destIp,
		DstPort: destPort,
		NetWork: network,
	}
	return metadata, nil
}
