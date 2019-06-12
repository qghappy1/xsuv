package module2

import (
	"fmt"
	"xsuv/nats"
	"xsuv/util/log"
	"encoding/binary"
	"xsuv/etcdv3/natslb"
)

const (
	ModuleMsgType = 0		// 无rpc模式
	ModuleRpcType = 1		// 无msg模式
	ModuleRpcMsgType = 2	// msg与rpc模式

	ModuleLbHash = 0		// 一致性策略
	ModuleLbRound = 1		// 循环策略
)

type IModule interface {
	OnInit()
	OnDestroy()
	Name() string
	Run(closeSig chan bool)
	OnMsgHandle(roleID int64, id uint16, msg []byte)
	OnRpcHandle(roleID int64, id uint16, msg []byte) []byte
	ConnectModule(moduleType string, lb int)
	SendMsg(moduleType string, roleID int64, msgID uint16, msg []byte) bool
	SendMsgSpecial(moduleName string, roleID int64, msgID uint16, msg []byte) bool
	RpcMsg(moduleType string, roleID int64, msgID uint16, msg []byte) []byte
}

//
type  BaseModule struct {
	name string
	rpcName string
	nat *nats.Nats
	c *natslb.Client
	lb *natslb.NatsLb
	otherModules map[string]natslb.IWatcher
	subModule IModule
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
		if err := m.nat.Subscribe(m.name, m.msgHandle); err != nil { log.FatalDepth(2, "%v", err) }
	case ModuleRpcType:
		if err := m.nat.Register(m.rpcName, m.rpcHandle); err != nil { log.FatalDepth(2, "%v", err) }
	case ModuleRpcMsgType:
		if err := m.nat.Subscribe(m.name, m.msgHandle); err != nil { log.FatalDepth(2, "%v", err) }
		if err := m.nat.Register(m.rpcName, m.rpcHandle); err != nil { log.FatalDepth(2, "%v", err) }
	default:
		log.FatalDepth(2, "msgType:%v err")
	}
	return m
}

func (m *BaseModule) OnInit() {}

func (m *BaseModule) OnMsgHandle(roleID int64, id uint16, msg []byte) { }

func (m *BaseModule) OnRpcHandle(roleID int64, id uint16, msg []byte) []byte { return nil }

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

func (m *BaseModule) SendMsg(moduleType string, roleID int64, msgID uint16, msg []byte) bool {
	if len(msg) == 0 { return false }
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
	data := make([]byte, 0)
	srole := make([]byte, 8)
	smsgid := make([]byte, 2)
	binary.BigEndian.PutUint64(srole, uint64(roleID))
	binary.BigEndian.PutUint16(smsgid, uint16(msgID))
	data = append(data, srole...)
	data = append(data, smsgid...)
	data = append(data, msg...)
	err = m.nat.Publish(name, data)
	if err != nil {
		log.ErrorDepth(2, "module:%v send msg.err:%v", m.name, err)
		return false
	}
	return true
}

func (m *BaseModule) SendMsgSpecial(moduleName string, roleID int64, msgID uint16, msg []byte) bool {
	if len(msg) == 0 { return false }
	data := make([]byte, 0)
	srole := make([]byte, 8)
	smsgid := make([]byte, 2)
	binary.BigEndian.PutUint64(srole, uint64(roleID))
	binary.BigEndian.PutUint16(smsgid, uint16(msgID))
	data = append(data, srole...)
	data = append(data, smsgid...)
	data = append(data, msg...)
	err := m.nat.Publish(moduleName, data)
	if err != nil {
		log.ErrorDepth(2, "module:%v send msg.err:%v", m.name, err)
		return false
	}
	return true
}

func (m *BaseModule) RpcMsg(moduleType string, roleID int64, msgID uint16, msg []byte) []byte {
	if len(msg) == 0 { return nil }
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
	data := make([]byte, 0)
	srole := make([]byte, 8)
	smsgid := make([]byte, 2)
	binary.BigEndian.PutUint64(srole, uint64(roleID))
	binary.BigEndian.PutUint16(smsgid, uint16(msgID))
	data = append(data, srole...)
	data = append(data, smsgid...)
	data = append(data, msg...)
	data, err = m.nat.Call(fmt.Sprint(name, "Rpc"), data, 30000)
	if err != nil {
		log.ErrorDepth(2, "module:%v rpc module:%v err:%v", m.name, name, err)
		return nil
	}
	return data
}

func (m *BaseModule) msgHandle(data []byte) []byte {
	defer log.ErrorPanic()
	if len(data)<10 { return []byte("") }
	roleID := int64(binary.BigEndian.Uint64(data))
	msgID := binary.BigEndian.Uint16(data[8:])
	if m.subModule != nil {
		m.subModule.OnMsgHandle(roleID, msgID, data[10:])
	}else{
		log.Error("module:%v handle nil", m.name)
	}
	return nil
}

func (m *BaseModule) rpcHandle(data []byte) []byte {
	defer log.ErrorPanic()
	if len(data)<10 { return []byte("") }
	roleID := int64(binary.BigEndian.Uint64(data))
	msgID := binary.BigEndian.Uint16(data[8:])
	if m.subModule != nil {
		return m.subModule.OnRpcHandle(roleID, msgID, data[10:])
	}else{
		log.Error("module:%v handle nil", m.name)
	}
	return []byte("")
}

