package netAgent

import (
	"fmt"

	"github.com/qghappy1/xsuv/network"
	"github.com/yuin/gopher-lua"
)

const (
	gateAgentTypeName = "agent"
)

func Preload(L *lua.LState) {
	registerGateAgentType(L)
	L.PreloadModule("gate", load)
}

func load(L *lua.LState) int {
	t := L.NewTable()
	L.Push(t)
	return 1
}

var luaGateAgentNameMethods = map[string]lua.LGFunction{
	"id":         gateAgentID,
	"sendmsg":    gateAgentSendMsg,
	"close":      gateAgentClose,
	"remoteaddr": gateAgentRemoteAddr,
}

func registerGateAgentType(L *lua.LState) {
	mt := L.NewTypeMetatable(gateAgentTypeName)
	L.SetGlobal(gateAgentTypeName, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), luaGateAgentNameMethods))
}

func GetLuaGateAgent(L *lua.LState, g network.GateAgent) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = g
	L.SetMetatable(ud, L.GetTypeMetatable(gateAgentTypeName))
	return ud
}

func checkGateAgent(L *lua.LState) network.GateAgent {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(network.GateAgent); ok {
		return v
	}
	L.ArgError(1, "not gate agent interface")
	return nil
}

func gateAgentID(L *lua.LState) int {
	g := checkGateAgent(L)
	L.Push(lua.LNumber(g.ID()))
	return 1
}

func gateAgentSendMsg(L *lua.LState) int {
	g := checkGateAgent(L)
	msg := L.CheckString(2)
	g.WriteMsg([]byte(msg))
	return 0
}

func gateAgentClose(L *lua.LState) int {
	g := checkGateAgent(L)
	g.Close()
	return 0
}

func gateAgentRemoteAddr(L *lua.LState) int {
	g := checkGateAgent(L)
	L.Push(lua.LString(g.RemoteAddr().String()))
	return 1
}
