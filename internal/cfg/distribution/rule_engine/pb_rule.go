package rule_engine

import "github.com/igoogolx/itun2socks/internal/constants"

func (x *PbRule) parse() (Rule, error) {

	return ParseItem(x.RuleType, x.Payload, x.Policy)

}

func (x *PbRule) Match(value string) bool {

	item, err := x.parse()

	if err != nil {
		return false
	}

	return item.Match(value)

}

func (x *PbRule) Value() string {

	item, err := x.parse()

	if err != nil {
		return ""
	}
	return item.Value()

}

func (x *PbRule) Type() constants.RuleType {

	item, err := x.parse()

	if err != nil {
		return ""
	}
	return item.Type()

}

func (x *PbRule) Valid() bool {

	item, err := x.parse()

	if err != nil {
		return false
	}
	return item.Valid()

}
