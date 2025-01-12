//go:build !darwin

package dns

func Hijack(_ string, server string, shouldReset bool) ([]string, error) {
	return []string{}, nil
}

func Resume(_ string, shouldReset bool) error {
	return nil
}
