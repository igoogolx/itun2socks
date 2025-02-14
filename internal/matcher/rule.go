package matcher

import (
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/ruleEngine"
	"sync"
)

var defaultRuleEngine *ruleEngine.Engine
var mux sync.RWMutex

func UpdateRuleEngine(e *ruleEngine.Engine) {
	mux.Lock()
	defer mux.Unlock()
	defaultRuleEngine = e
}

func GetRuleEngine() *ruleEngine.Engine {
	mux.RLock()
	defer mux.RUnlock()
	return defaultRuleEngine
}
