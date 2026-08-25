package meta_patch

import (
	clashC "github.com/igoogolx/itun2socks/pkg/clash/constant"
	metaC "github.com/metacubex/mihomo/constant"
	"net/netip"
)

func ConvertMeta(m *clashC.Metadata) *metaC.Metadata {

	srcIP, err := netip.ParseAddr(m.SrcIP.String())
	if err != nil {
		srcIP = netip.Addr{}
	}

	dstIP, err := netip.ParseAddr(m.DstIP.String())
	if err != nil {
		dstIP = netip.Addr{}
	}

	newMetadata := &metaC.Metadata{
		NetWork:     metaC.NetWork(m.NetWork),
		Type:        metaC.Type(m.Type),
		SrcIP:       srcIP,
		DstIP:       dstIP,
		SrcPort:     uint16(m.SrcPort),
		DstPort:     uint16(m.DstPort),
		Host:        m.Host,
		ProcessPath: m.ProcessPath,
	}

	return newMetadata

}
