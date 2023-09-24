//go:build !darwin

package dns

func Hijack(_ string) error {
	return nil
}

func Resume(_ string) error {
	return nil
}
