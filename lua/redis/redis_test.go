
package redis

import (
	"testing"
	"github.com/yuin/gopher-lua"
)


func TestLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	Preload(L)
	if err := L.DoFile("redis.lua"); err != nil {
		panic(err)
	}
}
