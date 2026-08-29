package dns

import (
	"github.com/metacubex/mihomo/component/fakeip"
	"net/netip"

	"github.com/igoogolx/itun2socks/internal/constants"
)

var ipRange, _ = netip.ParsePrefix(constants.TunGateway)

var FakeIpPool, _ = fakeip.New(fakeip.Options{
	IPNet: ipRange,
	Size:  4 * 1024,
})
