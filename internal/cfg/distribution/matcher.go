package distribution

type MatcherList interface {
	Has(s string) bool
}
