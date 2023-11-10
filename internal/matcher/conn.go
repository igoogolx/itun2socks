package matcher

import (
	C "github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/internal/constants"
	"sync"
)

type Conn interface {
	GetConnRule(metadata C.Metadata) constants.RuleType
}

var defaultConnMatcher Conn
var matcherMux sync.RWMutex

func UpdateConnMatcher(c Conn) {
	matcherMux.Lock()
	defer matcherMux.Unlock()
	defaultConnMatcher = c
}

func GetConnMatcher() Conn {
	matcherMux.RLock()
	defer matcherMux.RUnlock()
	return defaultConnMatcher
}
