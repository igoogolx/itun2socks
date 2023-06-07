package geo

import (
	C "github.com/Dreamacro/clash/constant"
	"os"
)

func init() {
	C.GeoIpUrl = "https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geoip.dat"
	C.GeoSiteUrl = "https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geosite.dat"
	curDir, _ := os.Getwd()
	C.GeodataMode = true
	C.SetHomeDir(curDir)
}
