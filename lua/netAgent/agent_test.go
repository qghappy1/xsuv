
package netAgent

import (
	"fmt"
	"net"
	"testing"
	"github.com/yuin/gopher-lua"
)



type TestAgent struct {
	userData interface{}
}

func (this *TestAgent) ID() int64{
	return 0
}

func (this *TestAgent) WriteMsg(msg []byte) {
	fmt.Printf("send msg:%v\n", msg)
}

func (this *TestAgent) LocalAddr() net.Addr{
	return new(net.TCPAddr)
}

func (this *TestAgent) RemoteAddr() net.Addr{
	return new(net.TCPAddr)
}

func (this *TestAgent) Close() {
	fmt.Printf("close\n")
}

func (this *TestAgent) Destroy() {
}

func (this *TestAgent) UserData() interface{}{
	return this.userData
}

func (this *TestAgent) SetUserData(data interface{}) {
	this.userData = data
}

func exampleGateAgent(L *lua.LState){
	agent := new(TestAgent)

	if err := L.DoFile("agent.lua"); err != nil {
		panic(err)
	}
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("test"),
		NRet:    0,
		Protect: true,
	}, GetLuaGateAgent(L, agent)); err != nil {
		panic(err)
	}
}

func TestLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	Preload(L)
	exampleGateAgent(L)
}
