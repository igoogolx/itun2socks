package rules

import (
	"errors"
	"slices"
)

var (
	errPayload = errors.New("payload error")

	noResolve = "no-resolve"
)

func HasNoResolve(params []string) bool {
	return slices.Contains(params, noResolve)
}
