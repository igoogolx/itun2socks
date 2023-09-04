package geo

import (
	"embed"
	"github.com/igoogolx/itun2socks/pkg/list"
	"github.com/igoogolx/itun2socks/pkg/log"
	"io/fs"
)

//go:embed geoData/*
var data embed.FS

func LoadGeoIPs(countries []string) ([]string, error) {
	ips := make([]string, 0)
	if len(countries) == 0 {
		return ips, nil
	}
	for _, country := range countries {
		items, err := parse("geoData/ip/" + country)
		if err != nil {
			return nil, err
		}
		ips = append(ips, items...)
	}
	return ips, nil
}

func parse(path string) ([]string, error) {
	file, err := data.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file fs.File) {
		err := file.Close()
		if err != nil {
			log.Warnln(log.FormatLog(log.ConfigurationPrefix, "fail to close geo file: %v"), path)
		}
	}(file)
	items, err := list.ParseFile(file)
	return items, nil
}
