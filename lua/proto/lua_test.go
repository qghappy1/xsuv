
package proto

import (
	"fmt"
	"testing"
	"xsuv/lua/luatest"
	"github.com/yuin/gopher-lua"
	proto2 "github.com/golang/protobuf/proto"
)

func TestLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	Preload(L)
	msg := new(luatest.CardBagSC)
	card := new(luatest.CardInfo)
	card.SetId(1)
	card.SetWeapon1(1)
	card.SetWeapon2(2)
	skill := new(luatest.CardSkill)
	skill.SetType(1)
	skill.SetLevel(2)
	skill.Id = []int32{1, 2, 3}
	card.Skills = append(card.Skills, skill)
	msg.Cards = append(msg.Cards, card)

	RegisterTypeToLua(L, "CardBagSC", luatest.CardBagSC{})
	RegisterTypeToLua(L, "CardInfo", luatest.CardInfo{})
	RegisterTypeToLua(L, "CardInfos", []*luatest.CardInfo{})

	data, err := proto2.Marshal(msg)
	if err != nil {
		panic(err)
	}
	RegisterMsgType(1, msg)
	RegisterMsgType(2, &luatest.CardUnlockSC{})
	luaMsg, err1 := MsgToLua(L, 1, data)
	if err1 != nil {
		panic(err1)
	}

	if err := L.DoFile("main.lua"); err != nil {
		panic(err)
	}
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("on_message"),
		NRet:    0,
		Protect: true,
	}, lua.LNumber(1), lua.LNumber(1), luaMsg); err != nil {
		panic(err)
	}
	//if err := L.CallByParam(lua.P{
	//	Fn:      L.GetGlobal("on_message"),
	//	NRet:    0,
	//	Protect: true,
	//}, lua.LNumber(1), lua.LNumber(1), RegisterValueToLua(L, msg)); err != nil {
	//	panic(err)
	//}
	fmt.Println("go skill.id", skill.Id)

	unlock := new(luatest.CardUnlockSC)
	unlock.Ret = luatest.ErrId_OK.Enum()
	unlock.SetId(10)
	unlock.SetSkillType(11)
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("on_message"),
		NRet:    0,
		Protect: true,
	}, lua.LNumber(2), lua.LNumber(2), RegisterValueToLua(L, unlock)); err != nil {
		panic(err)
	}
}
