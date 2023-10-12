package rule

type Rule interface {
	Match(value string) bool
	Value() string
	Policy() string
}
