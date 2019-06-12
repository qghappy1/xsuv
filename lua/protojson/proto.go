
package protojson

import (
	"fmt"
	"reflect"
	"encoding/json"
	"xsuv/util/log"
	"encoding/binary"
	json2 "xsuv/lua/json"
	"github.com/yuin/gopher-lua"
	proto2 "github.com/golang/protobuf/proto"
)

var (
	api = map[string]lua.LGFunction{
		"marshal": tableMarshalBytes,
		"unmarshal": unmarshalToTable,
	}
	msgInfo = make(map[uint16]reflect.Type)
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

func RegisterMsgType(msgID uint16, msg proto2.Message) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	msgInfo[msgID] = msgType
}

func _MarshalToJson(msg proto2.Message) string {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg:%v err:%v", err)
		return ""
	}
	return string(data)
}

func ProtoToJson(id uint16, data []byte) ([]byte, error) {
	msgType, ok := msgInfo[id]
	if !ok {
		return []byte(""), fmt.Errorf("msg:%v not register", id)
	}
	msg := reflect.New(msgType.Elem()).Interface()
	err := proto2.Unmarshal(data, msg.(proto2.Message))
	if err != nil {
		return []byte(""), err
	}
	data, err = json.Marshal(msg)
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}

func tableMarshalBytes(L *lua.LState) int{
	id := uint16(L.CheckInt(1))
	tbl := L.CheckTable(2)
	msgType, ok := msgInfo[id]
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("msg:%v not register", id)))
		return 2
	}
	msg := reflect.New(msgType.Elem()).Interface()
	bytes, err := json2.Encode(tbl)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("msg:%v lua table not encode json.err:%v", id, err.Error())))
		return 2
	}
	err = json.Unmarshal(bytes, msg)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("msg:%v lua table not unmarshal proto message.err:%v", id, err.Error())))
		return 2
	}

	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(id))
	buf, err := proto2.Marshal(msg.(proto2.Message))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("marshal msg:%v err:%v", id, err)))
		return 2
	}
	data = append(data, buf...)
	L.Push(lua.LString(data))
	return 1
}

func unmarshalToTable(L *lua.LState) int {
	str := L.CheckString(1)
	value, err := json2.Decode(L, []byte(str))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(value)
	return 1
}