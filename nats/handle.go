package nats

import (
	"xsuv/nats/api"
	"xsuv/util/log"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
)

type IHandle interface{
	Register(id uint16, f interface{})
	Handle(data []byte) []byte
}

//type HandleFunc struct {
//	functions map[uint16]func(string, int64, []byte)[]byte
//}
//
//func NewHandleFunc() *HandleFunc {
//	h := new(HandleFunc)
//	h.functions = make(map[uint16]func(string, int64, []byte)[]byte)
//	return h
//}
//
//func (this *HandleFunc) Register(id uint16, f func(string, int64, []byte)[]byte){
//	if _, ok := this.functions[id]; ok {
//		log.Fatal("function id:%v already registered", id)
//	}
//	this.functions[id] = f
//}
//
//func (this *HandleFunc) Handle(data []byte) []byte {
//	defer log.ErrorPanic()
//	iMsg := new(api.InnerMsg)
//	err := proto.Unmarshal(data, iMsg)
//	if err != nil {
//		log.Error("unmarshal err:%v", err)
//		return []byte("")
//	}
//	if len(iMsg.Msg)<2 {
//		log.Error("role:%v msg:%v from:%v error", iMsg.GetTokenId(), iMsg.Msg, iMsg.GetSrcServerName())
//		return []byte("")
//	}
//	id := binary.BigEndian.Uint16(iMsg.Msg)
//	if f, ok := this.functions[id]; ok {
//		return f(iMsg.GetSrcServerName(), iMsg.GetTokenId(), iMsg.Msg[2:])
//	}else{
//		log.Error("function id:%v not registered", id)
//		return []byte("")
//	}
//	return []byte("")
//}

// nats publish
func NatsPublish(pub *Nats, dstServer, srcServer string, roleID int64, data []byte) (error) {
	iMsg := new(api.InnerMsg)
	iMsg.SetSrcServerName(srcServer)
	iMsg.SetDstServerName(dstServer)
	iMsg.SetTokenId(roleID)
	iMsg.Msg = data
	data2, err := proto.Marshal(iMsg)
	if err != nil {
		log.Error("marshal err:%v", err)
		return err
	}
	if err = pub.Publish(dstServer, data2); err != nil {
		log.ErrorDepth(2,"call err:%v.dstSrv:%v", err, dstServer)
	}
	return err
}

// nats rpc调用
func NatsRpc(pub *Nats, dstServer, srcServer string, roleID int64, data []byte) ([]byte, error) {
	iMsg := new(api.InnerMsg)
	iMsg.SetSrcServerName(srcServer)
	iMsg.SetDstServerName(dstServer)
	iMsg.SetTokenId(roleID)
	iMsg.Msg = data
	data2, err := proto.Marshal(iMsg)
	if err != nil {
		log.Error("marshal err:%v", err)
		return nil, err
	}
	if ret, err := pub.Call(dstServer, data2, 30000); err != nil {
		log.ErrorDepth(2,"call err:%v.dstSrv:%v", err, dstServer)
		return nil, err
	}else{
		return ret, nil
	}
}

func Marshal(id uint16, msg proto.Message) []byte {
	sid := make([]byte, 2)
	binary.BigEndian.PutUint16(sid, id)
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("marshal msg:%v err:%v", id, err)
		return []byte("")
	}
	sid = append(sid, data...)
	return sid
}

func Unmarshal(data []byte, msg proto.Message) (error) {
	return proto.UnmarshalMerge(data, msg)
}