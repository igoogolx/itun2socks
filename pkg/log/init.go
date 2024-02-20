//go:build debug
// +build debug

package log

import cLog "github.com/igoogolx/clash/log"

func init() {

	cLog.SetLevel(cLog.DEBUG)

}
