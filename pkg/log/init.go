//go:build debug
// +build debug

package log

import cLog "github.com/Dreamacro/clash/log"

func main() {

	cLog.SetLevel(cLog.DEBUG)

}
