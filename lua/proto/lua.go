package proto

import (
	"encoding/binary"
	"fmt"
	"reflect"

	proto2 "github.com/golang/protobuf/proto"
	"github.com/layeh/gopher-luar"
	"github.com/qghappy1/xsuv/util/log"
	"github.com/yuin/gopher-lua"
)

var (
	RegisterValueToLua = luar.New
	msgInfo            = make(map[uint16]reflect.Type)
	api                = map[string]lua.LGFunction{
		"marshal": protoMarshal,
	}
)

func Preload(L *lua.LState) {
	L.PreloadModule("proto", load)
}

func load(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

// RegisterMsgType(1, &luatest.CardBagSC{})
func RegisterMsgType(msgID uint16, msg proto2.Message) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	msgInfo[msgID] = msgType
}

func MsgToLua(L *lua.LState, id uint16, buffer []byte) (lua.LValue, error) {
	msgType, ok := msgInfo[id]
	if !ok {
		return lua.LNil, fmt.Errorf("msg:% not register", id)
	}
	msg := reflect.New(msgType.Elem()).Interface()
	if err := proto2.Unmarshal(buffer, msg.(proto2.Message)); err != nil {
		return lua.LNil, err
	}
	return luar.New(L, msg), nil
}

func RegisterTypeToLua(L *lua.LState, tName string, v interface{}) {
	L.SetGlobal(tName, luar.NewType(L, v))
}

func protoMarshal(L *lua.LState) int {
	id := uint16(L.CheckInt(1))
	ud := L.CheckUserData(2)
	msg, ok := ud.Value.(proto2.Message)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("msg:%v not protobuf", reflect.TypeOf(ud.Value))))
		return 2
	}
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(id))
	buf, err := proto2.Marshal(msg)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("marshal msg:%v err:%v", id, err)))
		return 2
	}
	data = append(data, buf...)
	L.Push(lua.LString(data))
	return 1
}
