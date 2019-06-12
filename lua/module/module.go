package module

import (
	"github.com/qghappy1/xsuv/module2"
	"github.com/yuin/gopher-lua"
)

var luaModuleMethods = map[string]lua.LGFunction{
	"name":            moduleName,
	"sendmsg":         moduleSendMsg,
	"sendmsg_special": moduleSendMsgSpecial,
	"rpcmsg":          moduleRpcMsg,
}

func RegisterModule(L *lua.LState, moduleName string, m module2.IModule) {
	ud := L.NewUserData()
	ud.Value = m
	ud.Metatable = L.NewTypeMetatable(moduleName)
	L.SetField(ud.Metatable, "__index", L.SetFuncs(L.NewTable(), luaModuleMethods))
	L.SetGlobal(moduleName, ud)
}

func RegisterModuleType(L *lua.LState, moduleType, moduleTypeName string) {
	L.SetGlobal(moduleTypeName, lua.LString(moduleType))
}

func checkModule(L *lua.LState) module2.IModule {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(module2.IModule); ok {
		return v
	}
	L.ArgError(1, "not module interface")
	return nil
}

func moduleName(L *lua.LState) int {
	m := checkModule(L)
	L.Push(lua.LString(m.Name()))
	return 1
}

func moduleSendMsg(L *lua.LState) int {
	m := checkModule(L)
	moduleType := L.CheckString(2)
	roleID := L.CheckInt64(3)
	msgID := uint16(L.CheckInt(4))
	msg := L.CheckString(5)
	ok := m.SendMsg(moduleType, roleID, msgID, []byte(msg))
	L.Push(lua.LBool(ok))
	return 1
}

func moduleSendMsgSpecial(L *lua.LState) int {
	m := checkModule(L)
	moduleName := L.CheckString(2)
	roleID := L.CheckInt64(3)
	msgID := uint16(L.CheckInt(4))
	msg := L.CheckString(5)
	ok := m.SendMsgSpecial(moduleName, roleID, msgID, []byte(msg))
	L.Push(lua.LBool(ok))
	return 1
}

func moduleRpcMsg(L *lua.LState) int {
	m := checkModule(L)
	moduleType := L.CheckString(2)
	roleID := L.CheckInt64(3)
	msgID := uint16(L.CheckInt(4))
	msg := L.CheckString(5)
	ret := m.RpcMsg(moduleType, roleID, msgID, []byte(msg))
	L.Push(lua.LString(string(ret)))
	return 1
}
