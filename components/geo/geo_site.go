package geo

import (
	"github.com/igoogolx/itun2socks/components/list"
	"path/filepath"
)

func LoadGeoSites(rootDir string, countries []string) ([]string, error) {
	sites := make([]string, 0)
	if len(countries) == 0 {
		return sites, nil
	}
	for _, country := range countries {
		path := filepath.Join(rootDir, "site", country)
		items, err := list.ParseFile(path)
		if err != nil {
			return sites, nil
		}
		sites = append(sites, items...)
	}
	return sites, nil
}
