package tunnel

import (
	"errors"
	P "github.com/Dreamacro/clash/component/process"
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/log"
	"net"
	"net/netip"
	"strconv"
	"sync"
)

var defaultShouldFindProcess bool
var mux sync.RWMutex

func UpdateShouldFindProcess(value bool) {
	mux.Lock()
	defer mux.Unlock()
	defaultShouldFindProcess = value
}

func findProcessPath(metadata constant.Metadata) string {
	srcIP, ok := netip.AddrFromSlice(metadata.SrcIP)
	if ok && metadata.OriginDst.IsValid() {
		srcIP = srcIP.Unmap()
		path, err := P.FindProcessPath(metadata.NetWork.String(), netip.AddrPortFrom(srcIP, uint16(metadata.SrcPort)), metadata.OriginDst)
		if err != nil {
			log.Debugln("[Process] find process %s: %v", metadata.String(), err)
		} else {
			log.Debugln("[Process] %s from process %s", metadata.String(), path)
			return path
		}
	}
	return ""
}

func CreateUdpMetadata(srcAddr, destAddr net.UDPAddr) constant.Metadata {
	metadata := constant.Metadata{
		SrcIP:   srcAddr.IP,
		SrcPort: constant.Port(srcAddr.Port),
		DstIP:   destAddr.IP,
		DstPort: constant.Port(destAddr.Port),
		NetWork: constant.UDP,
	}
	if addrPort, err := netip.ParseAddrPort(destAddr.String()); err == nil {
		metadata.OriginDst = addrPort
	}

	if defaultShouldFindProcess {
		metadata.ProcessPath = findProcessPath(metadata)
	}
	return metadata
}

func CreateTcpMetadata(srcAddr, destAddr net.TCPAddr) constant.Metadata {
	metadata := constant.Metadata{
		SrcIP:   srcAddr.IP,
		SrcPort: constant.Port(srcAddr.Port),
		DstIP:   destAddr.IP,
		DstPort: constant.Port(destAddr.Port),
		NetWork: constant.TCP,
	}
	if addrPort, err := netip.ParseAddrPort(destAddr.String()); err == nil {
		metadata.OriginDst = addrPort
	}
	if defaultShouldFindProcess {
		metadata.ProcessPath = findProcessPath(metadata)
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
	metaSrcPort, err := strconv.Atoi(srcPort)
	if err != nil {
		return nil, err
	}
	metaDestPort, err := strconv.Atoi(destPort)

	if err != nil {
		return nil, err
	}

	metadata := &constant.Metadata{
		SrcIP:   srcIp,
		SrcPort: constant.Port(metaSrcPort),
		DstIP:   destIp,
		DstPort: constant.Port(metaDestPort),
		NetWork: network,
	}
	return metadata, nil
}
