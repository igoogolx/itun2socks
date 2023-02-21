package global

//TODO: remove global module
import "go.uber.org/atomic"

var defaultInterfaceName = atomic.NewString("")

func UpdateDefaultInterfaceName(name string) {
	defaultInterfaceName.Store(name)
}

func GetDefaultInterfaceName() string {
	return defaultInterfaceName.Load()
}
