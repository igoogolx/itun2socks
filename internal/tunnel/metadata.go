package tunnel

import (
	"errors"
	P "github.com/igoogolx/itun2socks/pkg/clash/component/process"
	"github.com/igoogolx/itun2socks/pkg/log"
	"net"
	"net/netip"
	"strconv"
	"sync"

	metaC "github.com/metacubex/mihomo/constant"
)

var defaultShouldFindProcess bool
var mux sync.RWMutex

func UpdateShouldFindProcess(value bool) {
	mux.Lock()
	defer mux.Unlock()
	defaultShouldFindProcess = value
}

func findProcessPath(metadata metaC.Metadata) string {

	path, err := P.FindProcessPath(metadata.NetWork.String(), netip.AddrPortFrom(metadata.SrcIP, metadata.SrcPort), netip.AddrPortFrom(metadata.DstIP, metadata.DstPort))
	if err != nil {
		log.Debugln(log.FormatLog(log.RulePrefix, "find process %s: %v"), metadata.String(), err)
	} else {
		log.Debugln(log.FormatLog(log.RulePrefix, "find process %s: %v"), metadata.String(), path)
		return path
	}
	return ""
}

func CreateUdpMetadata(srcAddr, destAddr netip.AddrPort) metaC.Metadata {
	metadata := metaC.Metadata{
		SrcIP:   srcAddr.Addr(),
		SrcPort: srcAddr.Port(),
		DstIP:   destAddr.Addr(),
		DstPort: destAddr.Port(),
		NetWork: metaC.UDP,
	}

	if defaultShouldFindProcess {
		metadata.ProcessPath = findProcessPath(metadata)
	}
	return metadata
}

func CreateTcpMetadata(srcAddr, destAddr netip.AddrPort) metaC.Metadata {
	metadata := metaC.Metadata{
		SrcIP:   srcAddr.Addr(),
		SrcPort: srcAddr.Port(),
		DstIP:   destAddr.Addr(),
		DstPort: destAddr.Port(),
		NetWork: metaC.TCP,
	}
	if defaultShouldFindProcess {
		metadata.ProcessPath = findProcessPath(metadata)
	}
	return metadata
}

func CreateMetadata(srcAddr, destAddr string, network metaC.NetWork) (*metaC.Metadata, error) {
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

	metadata := &metaC.Metadata{
		SrcIP:   srcIp,
		SrcPort: uint16(metaSrcPort),
		DstIP:   destIp,
		DstPort: uint16(metaDestPort),
		NetWork: network,
	}
	return metadata, nil
}
