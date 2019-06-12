
package protojson

import (
	"testing"
	"xsuv/lua/luatest"
	"github.com/yuin/gopher-lua"
	json2 "github.com/layeh/gopher-json"
)


func TestLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	Preload(L)
	json2.Preload(L)

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

	RegisterMsgType(1, &luatest.CardBagSC{})
	RegisterMsgType(2, &luatest.CardUnlockSC{})

	s := _MarshalToJson(msg)
	if err := L.DoFile("main.lua"); err != nil {
		panic(err)
	}
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("on_message"),
		NRet:    0,
		Protect: true,
	}, lua.LNumber(1), lua.LNumber(1), lua.LString(s)); err != nil {
		panic(err)
	}

	unlock := new(luatest.CardUnlockSC)
	unlock.Ret = luatest.ErrId_OK.Enum()
	unlock.SetId(10)
	unlock.SetSkillType(11)
	s = _MarshalToJson(unlock)
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("on_message"),
		NRet:    0,
		Protect: true,
	}, lua.LNumber(2), lua.LNumber(2), lua.LString(s)); err != nil {
		panic(err)
	}
}
