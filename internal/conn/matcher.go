package conn

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"sync"
)

type Matcher interface {
	GetRule(metadata C.Metadata) constants.RuleType
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
