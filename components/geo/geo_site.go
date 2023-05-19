package geo

func LoadGeoSites(countries []string) ([]string, error) {
	sites := make([]string, 0)
	if len(countries) == 0 {
		return sites, nil
	}
	for _, country := range countries {
		items, err := parse("geoData/site/" + country)
		if err != nil {
			return nil, err
		}
		sites = append(sites, items...)
	}
	return sites, nil
}
