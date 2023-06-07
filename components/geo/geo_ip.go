package geo

import (
	"github.com/Dreamacro/clash/component/geodata"
	"github.com/Dreamacro/clash/component/geodata/router"
)

func LoadGeoIPs(countries []string) ([]*router.GeoIPMatcher, error) {
	err := geodata.InitGeoIP()
	if err != nil {
		return nil, err
	}
	ips := make([]*router.GeoIPMatcher, 0)
	if len(countries) == 0 {
		return ips, nil
	}
	for _, country := range countries {
		items, _, err := geodata.LoadGeoIPMatcher(country)
		if err != nil {
			return nil, err
		}
		ips = append(ips, items)
	}
	return ips, nil
}
