package conn

import (
	"sync"
)

var defaultConnMatchers []Matcher
var matcherMux sync.RWMutex

func UpdateConnMatcher(matchers []Matcher) {
	matcherMux.Lock()
	defer matcherMux.Unlock()
	defaultConnMatchers = matchers
}

func GetConnMatcher() []Matcher {
	matcherMux.RLock()
	defer matcherMux.RUnlock()
	return defaultConnMatchers
}
