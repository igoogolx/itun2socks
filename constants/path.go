package constants

import (
	"path/filepath"
)

var Path path

type path struct {
	homeDir string
}

func (p path) SetHomeDir(dir string) {
	p.homeDir = dir
}

func (p path) HomeDir() string {
	return p.homeDir
}

func (p path) GeoDataDir() string {
	return filepath.Join(HomeDir, "geoData")
}

func (p path) WebDir() string {
	return filepath.Join(HomeDir, "web", "dist")
}

func (p path) ConfigFilePath() string {
	return filepath.Join(HomeDir, DbFileName)
}
