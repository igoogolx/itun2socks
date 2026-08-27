package tunnel

import (
	"errors"
	P "github.com/metacubex/mihomo/component/process"
	"net"
	"net/netip"
	"strconv"
	"sync"

	"github.com/igoogolx/itun2socks/pkg/clash/constant"
	"github.com/igoogolx/itun2socks/pkg/log"
)

var defaultShouldFindProcess bool
var mux sync.RWMutex

func UpdateShouldFindProcess(value bool) {
	mux.Lock()
	defer mux.Unlock()
	defaultShouldFindProcess = value
}

func findProcessPath(metadata constant.Metadata) string {
	if metadata.OriginDst.IsValid() {
		_, path, err := P.FindProcessName(metadata.NetWork.String(), metadata.SrcIP, int(metadata.SrcPort))
		if err != nil {
			log.Debugln(log.FormatLog(log.RulePrefix, "find process %s: %v"), metadata.String(), err)
		} else {
			log.Debugln(log.FormatLog(log.RulePrefix, "find process %s: %v"), metadata.String(), path)
			return path
		}
	}
	return ""
}

func CreateUdpMetadata(srcAddr, destAddr netip.AddrPort) constant.Metadata {
	metadata := constant.Metadata{
		SrcIP:   srcAddr.Addr(),
		SrcPort: constant.Port(srcAddr.Port()),
		DstIP:   destAddr.Addr(),
		DstPort: constant.Port(destAddr.Port()),
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

func CreateTcpMetadata(srcAddr, destAddr netip.AddrPort) constant.Metadata {
	metadata := constant.Metadata{
		SrcIP:   srcAddr.Addr(),
		SrcPort: constant.Port(srcAddr.Port()),
		DstIP:   destAddr.Addr(),
		DstPort: constant.Port(destAddr.Port()),
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
