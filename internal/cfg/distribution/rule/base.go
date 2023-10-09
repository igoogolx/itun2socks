package rule

type Rule interface {
	Match(value string) bool
}
