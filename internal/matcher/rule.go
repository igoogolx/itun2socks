package matcher

import (
	"github.com/igoogolx/itun2socks/internal/cfg/distribution/rule_engine"
	"sync"
)

var defaultRuleEngine *rule_engine.Engine
var mux sync.RWMutex

func UpdateRuleEngine(e *rule_engine.Engine) {
	mux.Lock()
	defer mux.Unlock()
	defaultRuleEngine = e
}

func GetRuleEngine() *rule_engine.Engine {
	mux.RLock()
	defer mux.RUnlock()
	return defaultRuleEngine
}
