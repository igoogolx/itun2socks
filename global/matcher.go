package global

import (
	"github.com/igoogolx/itun2socks/constants"
	"sync"
)

type Matcher interface {
	GetRule(ip string) constants.IpRule
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
