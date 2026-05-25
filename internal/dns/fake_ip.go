package dns

import (
	"net"

	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/igoogolx/itun2socks/pkg/clash/component/fakeip"
)

var _, ipRange, _ = net.ParseCIDR(constants.FakeIpRange)

var FakeIpPool, _ = fakeip.New(fakeip.Options{
	IPNet: ipRange,
	Host:  nil,
	Size:  4 * 1024,
})
