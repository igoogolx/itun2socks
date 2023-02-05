package outbound

import (
	"fmt"
	"github.com/Dreamacro/clash/constant"
	"github.com/igoogolx/itun2socks/parser"
)

func New(proxy []map[string]interface{}, selected string) (constant.Proxy, error) {
	var selectedProxy map[string]interface{}
	for _, v := range proxy {
		if v["id"] == selected {
			selectedProxy = v
			break
		}
	}
	if selectedProxy == nil {
		return nil, fmt.Errorf("error getting seleted proxy, id:%v", selected)
	}
	p, err := parser.ParseProxy(selectedProxy)
	if err != nil {
		return nil, err
	}
	return p, nil
}
