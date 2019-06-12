
package module

import (
	"fmt"
	"time"
	"testing"
	"xsuv/nats"
	"xsuv/module2"
	"xsuv/etcdv3/natslb"
	"github.com/yuin/gopher-lua"
)



type gameModuleExample struct {
	*module2.BaseModule
}

type loginModuleExample struct {
	*module2.BaseModule
}

func (m *gameModuleExample) OnHandle(src string, roleID int64, id uint16, msg []byte) []byte {
	fmt.Printf("%v recv <%v> from %v.time:%v\n", m.Name(), string(msg), src, time.Now().Unix())
	return []byte(fmt.Sprint("i am ", m.Name()))
}

func (m *gameModuleExample) OnRpcHandle(src string, roleID int64, id uint16, msg []byte) []byte {
	fmt.Printf("%v recv <%v> from %v.time:%v\n", m.Name(), string(msg), src, time.Now().Unix())
	return []byte(fmt.Sprint("i am ", m.Name()))
}

func (m *loginModuleExample) OnHandle(src string, roleID int64, id uint16, msg []byte) []byte {
	fmt.Printf("%v recv <%v> from %v.time:%v\n", m.Name(), string(msg), src, time.Now().Unix())
	return []byte(fmt.Sprint("i am ", m.Name()))
}

func exampleModule(L *lua.LState){
	natUrls := "nats://111.230.46.154:4242"
	natUsername := "ft2018"
	natPassword := "T0pS3cr3t"
	nat := nats.NewNats(natUrls, natUsername, natPassword)
	etcAddress := "http://111.230.46.154:2379"
	c := natslb.ConnectEtcd(etcAddress)

	gameType := "TestGame"
	gameName := "TestGame1"
	g := new(gameModuleExample)
	g.BaseModule = module2.NewBaseModule(gameType, gameName, module2.ModuleRpcMsgType, nat, c, g)

	loginType := "TestLogin"
	loginName := "TestLogin1"
	l := new(loginModuleExample)
	l.BaseModule = module2.NewBaseModule(loginType, loginName, module2.ModuleRpcMsgType, nat, c, l)

	g.ConnectModule(loginType, module2.ModuleLbHash)
	l.ConnectModule(gameType, module2.ModuleLbHash)

	RegisterModule(L, gameName, g)
	RegisterModule(L, loginName, l)
	RegisterModuleType(L, gameType, "GameType")
	RegisterModuleType(L, loginType, "LoginType")
	if err := L.DoFile("module.lua"); err != nil {
		panic(err)
	}
	//g.SendMsg(loginType, 0, []byte(fmt.Sprintf("  %v", time.Now().Unix())))
	//msg := l.RpcMsg(gameType, 0, []byte(fmt.Sprintf("  %v", time.Now().Unix())))
	//fmt.Printf("%v rpc <%v> from %v.time:%v\n", l.Name(), string(msg), g.Name(), time.Now().Unix())
}

func TestLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	exampleModule(L)
}
