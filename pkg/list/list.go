package list

import (
	"bufio"
	"io"
	"strings"
)

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
			break
		}
	}
	return items, nil
}
