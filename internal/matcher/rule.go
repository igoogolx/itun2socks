package matcher

import (
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"sync"
)

var defaultRule *ruleEngine.Engine
var mux sync.RWMutex

func UpdateRule(e *ruleEngine.Engine) {
	mux.Lock()
	defer mux.Unlock()
	defaultRule = e
}

func GetRule() *ruleEngine.Engine {
	mux.RLock()
	defer mux.RUnlock()
	return defaultRule
}
