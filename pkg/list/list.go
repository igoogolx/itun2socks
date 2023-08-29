package list

import (
	"bufio"
	"github.com/Dreamacro/clash/log"
	"io"
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

func ParseFile(file io.Reader) ([]string, error) {
	items := make([]string, 0)
	var err error
	scanner := bufio.NewScanner(file)
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
