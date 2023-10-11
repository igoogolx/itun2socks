package conn

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"sync"
)

type Matcher interface {
	GetRule(ip string) constants.RuleType
}

var defaultMatcher Matcher
var matcherMux sync.RWMutex

func UpdateMatcher(c Matcher) {
	matcherMux.Lock()
	defer matcherMux.Unlock()
	defaultMatcher = c
}

func GetMatcher() Matcher {
	matcherMux.RLock()
	defer matcherMux.RUnlock()
	return defaultMatcher
}
