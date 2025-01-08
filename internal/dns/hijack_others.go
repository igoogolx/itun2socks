//go:build !darwin

package dns

func Hijack(_ string, server string) ([]string, error) {
	return []string{}, nil
}

func Resume(_ string) error {
	return nil
}
