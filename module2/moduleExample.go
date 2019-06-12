package module2

import (
	"fmt"
	"time"

	"github.com/qghappy1/xsuv/etcdv3/natslb"
	"github.com/qghappy1/xsuv/nats"
)

type gameModuleExample struct {
	*BaseModule
}

type loginModuleExample struct {
	*BaseModule
}

func (m *gameModuleExample) OnHandle(roleID int64, id uint16, msg []byte) []byte {
	fmt.Printf("%v recv <%v>.time:%v\n", m.Name(), string(msg), time.Now().Unix())
	return []byte(fmt.Sprint("i am ", m.Name()))
}

func (m *loginModuleExample) OnHandle(roleID int64, id uint16, msg []byte) []byte {
	fmt.Printf("%v recv <%v>.time:%v\n", m.Name(), string(msg), time.Now().Unix())
	return []byte(fmt.Sprint("i am ", m.Name()))
}

func exampleModule() {
	natUrls := "nats://111.230.46.154:4242"
	natUsername := "ft2018"
	natPassword := "T0pS3cr3t"
	nat := nats.NewNats(natUrls, natUsername, natPassword)
	etcAddress := "http://111.230.46.154:2379"
	c := natslb.ConnectEtcd(etcAddress)

	gameType := "TestGame"
	gameName := "TestGame1"
	g := new(gameModuleExample)
	g.BaseModule = NewBaseModule(gameType, gameName, ModuleRpcMsgType, nat, c, g)

	loginType := "TestLogin"
	loginName := "TestLogin1"
	l := new(loginModuleExample)
	l.BaseModule = NewBaseModule(loginType, loginName, ModuleRpcMsgType, nat, c, l)

	g.ConnectModule(loginType, ModuleLbHash)
	l.ConnectModule(gameType, ModuleLbHash)

	g.SendMsg(loginType, 0, 0, []byte(fmt.Sprintf("  %v", time.Now().Unix())))
	msg := l.RpcMsg(gameType, 0, 0, []byte(fmt.Sprintf("  %v", time.Now().Unix())))
	fmt.Printf("%v rpc <%v> from %v.time:%v\n", l.Name(), string(msg), g.Name(), time.Now().Unix())
}
