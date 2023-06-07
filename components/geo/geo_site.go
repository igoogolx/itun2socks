package geo

import (
	"github.com/Dreamacro/clash/component/geodata"
	"github.com/Dreamacro/clash/component/geodata/router"
)

func LoadGeoSites(countries []string) ([]*router.DomainMatcher, error) {
	err := geodata.InitGeoSite()
	if err != nil {
		return nil, err
	}

	sites := make([]*router.DomainMatcher, 0)
	if len(countries) == 0 {
		return sites, nil
	}
	for _, country := range countries {
		items, _, err := geodata.LoadGeoSiteMatcher(country)
		if err != nil {
			return nil, err
		}
		sites = append(sites, items)
	}
	return sites, nil
}
