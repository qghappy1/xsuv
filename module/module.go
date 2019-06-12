package module

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/qghappy1/xsuv/etcdv3/natslb"
	"github.com/qghappy1/xsuv/nats"
	"github.com/qghappy1/xsuv/nats/api"
	"github.com/qghappy1/xsuv/util/log"
)

const (
	ModuleMsgType    = 0 // 无rpc模式
	ModuleRpcType    = 1 // 无msg模式
	ModuleRpcMsgType = 2 // msg与rpc模式

	ModuleLbHash  = 0 // 一致性策略
	ModuleLbRound = 1 // 循环策略
)

type IModule interface {
	OnInit()
	OnDestroy()
	Name() string
	Run(closeSig chan bool)
	OnMsgHandle(src string, roleID int64, id uint16, msg []byte)
	OnRpcHandle(src string, roleID int64, id uint16, msg []byte) []byte
	ConnectModule(moduleType string, lb int)
	SendMsg(moduleType string, roleID int64, msg []byte) bool
	SendMsgSpecial(moduleName string, roleID int64, msg []byte) bool
	RpcMsg(moduleType string, roleID int64, msg []byte) []byte
}

//
type BaseModule struct {
	name         string
	rpcName      string
	nat          *nats.Nats
	c            *natslb.Client
	lb           *natslb.NatsLb
	otherModules map[string]natslb.IWatcher
	subModule    IModule
}

// subModule主要是获取子类的OnHandle实现
func NewBaseModule(moduleType string, moduleName string, msgType int, nat *nats.Nats, c *natslb.Client, subModule IModule) *BaseModule {
	m := new(BaseModule)
	m.name = moduleName
	m.rpcName = fmt.Sprint(m.name, "Rpc")
	m.nat = nat
	m.c = c
	m.otherModules = make(map[string]natslb.IWatcher)
	m.lb = natslb.NewNatsLb(c)
	m.lb.Register(moduleType, moduleName)
	m.subModule = subModule
	switch msgType {
	case ModuleMsgType:
		if err := m.nat.Subscribe(m.name, m.msgHandle); err != nil {
			log.FatalDepth(2, "%v", err)
		}
	case ModuleRpcType:
		if err := m.nat.Register(m.rpcName, m.rpcHandle); err != nil {
			log.FatalDepth(2, "%v", err)
		}
	case ModuleRpcMsgType:
		if err := m.nat.Subscribe(m.name, m.msgHandle); err != nil {
			log.FatalDepth(2, "%v", err)
		}
		if err := m.nat.Register(m.rpcName, m.rpcHandle); err != nil {
			log.FatalDepth(2, "%v", err)
		}
	default:
		log.FatalDepth(2, "msgType:%v err")
	}
	return m
}

func (m *BaseModule) OnInit() {}

func (m *BaseModule) OnMsgHandle(src string, roleID int64, id uint16, msg []byte) {}

func (m *BaseModule) OnRpcHandle(src string, roleID int64, id uint16, msg []byte) []byte { return nil }

func (m *BaseModule) ConnectModule(moduleType string, lb int) {
	switch lb {
	case ModuleLbHash:
		w := natslb.NewWatcherKetama(moduleType, m.c)
		w.Start()
		m.otherModules[moduleType] = w
	case ModuleLbRound:
		w := natslb.NewWatcherRound(moduleType, m.c)
		w.Start()
		m.otherModules[moduleType] = w
	default:
		log.FatalDepth(2, "lb:%v err")
	}
}

func (m *BaseModule) OnDestroy() { m.lb.UnRegister() }

func (m *BaseModule) Name() string { return m.name }

func (m *BaseModule) Run(closeSig chan bool) {}

func (m *BaseModule) SendMsg(moduleType string, roleID int64, msg []byte) bool {
	if len(msg) == 0 {
		return false
	}
	w, ok := m.otherModules[moduleType]
	if !ok {
		log.ErrorDepth(2, "not connect module:%v", moduleType)
		return false
	}
	name, err := w.Get(moduleType)
	if err != nil {
		log.ErrorDepth(2, "module:%v not register etcd", moduleType)
		return false
	}
	iMsg := new(api.InnerMsg)
	iMsg.SetSrcServerName(m.name)
	iMsg.SetDstServerName(name)
	iMsg.SetTokenId(roleID)
	iMsg.Msg = msg
	data, err1 := proto.Marshal(iMsg)
	if err1 != nil {
		log.ErrorDepth(2, "module:%v marshal err:%v", m.name, err1)
		return false
	}
	err = m.nat.Publish(name, data)
	if err != nil {
		log.ErrorDepth(2, "module:%v send msg.err:%v", m.name, err)
		return false
	}
	return true
}

func (m *BaseModule) SendMsgSpecial(moduleName string, roleID int64, msg []byte) bool {
	if len(msg) == 0 {
		return false
	}
	iMsg := new(api.InnerMsg)
	iMsg.SetSrcServerName(m.name)
	iMsg.SetDstServerName(moduleName)
	iMsg.SetTokenId(roleID)
	iMsg.Msg = msg
	data, err := proto.Marshal(iMsg)
	if err != nil {
		log.ErrorDepth(2, "module:%v marshal err:%v", m.name, err)
		return false
	}
	err = m.nat.Publish(moduleName, data)
	if err != nil {
		log.ErrorDepth(2, "module:%v send msg.err:%v", m.name, err)
		return false
	}
	return true
}

func (m *BaseModule) RpcMsg(moduleType string, roleID int64, msg []byte) []byte {
	if len(msg) == 0 {
		return nil
	}
	w, ok := m.otherModules[moduleType]
	if !ok {
		log.ErrorDepth(2, "not connect module:%v", moduleType)
		return nil
	}
	name, err := w.Get(moduleType)
	if err != nil {
		log.ErrorDepth(2, "module:%v not register etcd", moduleType)
		return nil
	}
	iMsg := new(api.InnerMsg)
	iMsg.SetSrcServerName(m.name)
	iMsg.SetDstServerName(name)
	iMsg.SetTokenId(roleID)
	iMsg.Msg = msg
	data, err1 := proto.Marshal(iMsg)
	if err1 != nil {
		log.ErrorDepth(2, "module:%v marshal err:%v", m.name, err1)
		return nil
	}
	retMsg, err2 := m.nat.Call(fmt.Sprint(name, "Rpc"), data, 30000)
	if err2 != nil {
		log.ErrorDepth(2, "module:%v rpc module:%v err:%v", m.name, name, err2)
		return nil
	}
	if retMsg == nil || len(retMsg) < 2 {
		return nil
	}
	return retMsg
}

func (m *BaseModule) msgHandle(data []byte) []byte {
	defer log.ErrorPanic()
	iMsg := new(api.InnerMsg)
	if err := proto.Unmarshal(data, iMsg); err != nil {
		log.Error("module:%v unmarshal err:%v", m.name, err)
		return nil
	}
	roleID := iMsg.GetTokenId()
	if len(iMsg.Msg) < 2 {
		log.Debug("module:%v role:%v msg less", m.name, roleID)
		return nil
	}
	id := binary.BigEndian.Uint16(iMsg.Msg)
	if m.subModule != nil {
		m.subModule.OnMsgHandle(iMsg.GetSrcServerName(), roleID, id, iMsg.Msg)
	} else {
		log.Error("module:%v handle nil", m.name)
	}
	return nil
}

func (m *BaseModule) rpcHandle(data []byte) []byte {
	defer log.ErrorPanic()
	iMsg := new(api.InnerMsg)
	if err := proto.Unmarshal(data, iMsg); err != nil {
		log.Error("module:%v unmarshal err:%v", m.name, err)
		return []byte("")
	}
	roleID := iMsg.GetTokenId()
	if len(iMsg.Msg) < 2 {
		log.Debug("module:%v role:%v msg less", m.name, roleID)
		return nil
	}
	id := binary.BigEndian.Uint16(iMsg.Msg)
	if m.subModule != nil {
		return m.subModule.OnRpcHandle(iMsg.GetSrcServerName(), roleID, id, iMsg.Msg)
	} else {
		log.Error("module:%v handle nil", m.name)
	}
	return []byte("")
}
