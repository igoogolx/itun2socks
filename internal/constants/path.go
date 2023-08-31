package constants

import (
	"path/filepath"
)

var Path path

type path struct {
	homeDir string
}

func (p *path) SetHomeDir(dir string) {
	p.homeDir = dir
}

func (p *path) HomeDir() string {
	return p.homeDir
}

func (p *path) ConfigFilePath() string {
	return filepath.Join(Path.HomeDir(), DbFileName)
}

func (p *path) LogFilePath() string {
	return filepath.Join(Path.HomeDir(), "core.log")
}
