package list

import (
	"bufio"
	"fmt"
	"github.com/Dreamacro/clash/log"
	"io"
	"os"
	"strings"
)

type Lister struct {
	Items  []string
	Mather func(s, i string) bool
}

func (l *Lister) Has(i string) bool {
	for _, item := range l.Items {
		if l.Mather(item, i) {
			return true
		}
	}
	return false
}

func (l *Lister) Insert(s string) error {
	l.Items = append(l.Items, s)
	return nil
}

func ParseFile(file string) ([]string, error) {
	items := make([]string, 0)
	if len(file) == 0 {
		return items, nil
	}
	f, err := os.Open(file)
	if err != nil {
		return items, fmt.Errorf("fail to open file %s: %s", file, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Warnln("fail to close file, file: %v, err: %v", file, err)
		}
	}(f)
	scanner := bufio.NewScanner(f)
	s := strings.Builder{}
	for scanner.Scan() {
		s.Write(scanner.Bytes())
		line := s.String()
		s.Reset()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(line)
		if line != "" {
			items = append(items, line)
		}
		if line == "" && err == io.EOF {
			log.Debugln("Reading file %s reached EOF", file)
			break
		}
	}
	return items, nil
}

func New(items []string, matcher func(s, i string) bool) *Lister {
	return &Lister{
		Items:  items,
		Mather: matcher,
	}
}
