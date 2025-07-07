package dns

import (
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net"
)

var _, ipRange, _ = net.ParseCIDR(constants.FakeIpRange)

var FakeIpPool, _ = fakeip.New(fakeip.Options{
	IPNet: ipRange,
	Host:  nil,
	Size:  4 * 1024,
})
