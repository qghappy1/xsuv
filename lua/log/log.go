package log

import (
	log2 "github.com/qghappy1/xsuv/util/log"
	"github.com/yuin/gopher-lua"
)

func logDebug(L *lua.LState) int {
	s := L.CheckString(1)
	log2.LuaDebug(s)
	return 0
}

func logInfo(L *lua.LState) int {
	s := L.CheckString(1)
	log2.LuaInfo(s)
	return 0
}

func logError(L *lua.LState) int {
	s := L.CheckString(1)
	log2.LuaError(s)
	return 0
}

func RegisterLog(L *lua.LState) {
	L.SetGlobal("LogDebug", L.NewFunction(logDebug))
	L.SetGlobal("LogInfo", L.NewFunction(logInfo))
	L.SetGlobal("LogError", L.NewFunction(logError))
}
