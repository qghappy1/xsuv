
package log

import (
	"testing"
	"github.com/yuin/gopher-lua"
)

func TestLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	RegisterLog(L)
	if err := L.DoFile("log.lua"); err != nil {
		panic(err)
	}
}
