package parser

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/constant"
)

func ParseProxy(mapping map[string]interface{}) (constant.Proxy, error) {
	return adapter.ParseProxy(mapping)
}
