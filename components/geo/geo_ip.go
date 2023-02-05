package geo

import (
	"github.com/igoogolx/itun2socks/components/list"
	"path/filepath"
)

func LoadGeoIPs(countries []string) ([]string, error) {
	ips := make([]string, 0)
	if len(countries) == 0 {
		return ips, nil
	}
	for _, country := range countries {
		path := filepath.Join("geoData", "ip", country)
		items, err := list.ParseFile(path)
		if err != nil {
			return ips, nil
		}
		ips = append(ips, items...)
	}
	return ips, nil
}
