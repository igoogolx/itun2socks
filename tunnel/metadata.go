package tunnel

import (
	"errors"
	"github.com/Dreamacro/clash/constant"
	"net"
	"net/netip"
	"strconv"
)

func CreateUdpMetadata(srcAddr, destAddr net.UDPAddr) (constant.Metadata, error) {
	srcIp, err := netip.ParseAddr(srcAddr.IP.String())
	if err != nil {
		return constant.Metadata{}, err
	}
	dstIp, err := netip.ParseAddr(destAddr.IP.String())
	if err != nil {
		return constant.Metadata{}, err
	}
	metadata := constant.Metadata{
		SrcIP:   srcIp,
		SrcPort: strconv.Itoa(srcAddr.Port),
		DstIP:   dstIp,
		DstPort: strconv.Itoa(destAddr.Port),
		NetWork: constant.UDP,
	}
	return metadata, nil
}

func CreateTcpMetadata(srcAddr, destAddr net.TCPAddr) (constant.Metadata, error) {
	srcIp, err := netip.ParseAddr(srcAddr.IP.String())
	if err != nil {
		return constant.Metadata{}, err
	}
	dstIp, err := netip.ParseAddr(destAddr.IP.String())
	if err != nil {
		return constant.Metadata{}, err
	}
	metadata := constant.Metadata{
		SrcIP:   srcIp,
		SrcPort: strconv.Itoa(srcAddr.Port),
		DstIP:   dstIp,
		DstPort: strconv.Itoa(destAddr.Port),
		NetWork: constant.TCP,
	}
	return metadata, nil
}

func CreateMetadata(srcAddr, destAddr string, network constant.NetWork) (*constant.Metadata, error) {
	var srcHost, srcPort string
	var srcIp netip.Addr
	var err error
	if len(srcAddr) != 0 {
		srcHost, srcPort, err = net.SplitHostPort(srcAddr)
		if err != nil {
			return nil, err
		}
		srcIp, err = netip.ParseAddr(srcHost)
		if err != nil {
			return nil, errors.New("fail to parse src host")
		}
	}

	destHost, destPort, err := net.SplitHostPort(destAddr)

	destIp, err := netip.ParseAddr(destHost)
	if err != nil {
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
