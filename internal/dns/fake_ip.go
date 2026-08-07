package dns

import (
	"net/netip"

	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/clash/component/fakeip"
)

var ipRange, _ = netip.ParsePrefix(constants.FakeIpRange)

var FakeIpPool, _ = fakeip.New(fakeip.Options{
	IPNet: &ipRange,
	Host:  nil,
	Size:  4 * 1024,
})
