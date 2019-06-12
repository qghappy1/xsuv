package network

import (
	"net"
	"xsuv/util/log"
)

var next = int64(1)

type Gate struct {
	MaxConnNum		int
	PendingWriteNum	int
	MaxMsgLen		uint32
	Handle			func([]byte, GateAgent) bool
	OnCloseAgent	func(GateAgent)
	OnNewAgent		func(GateAgent)

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

func (gate *Gate) Run(closeSig chan bool) {
	var tcpServer *TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *TCPConn) Agent {
			a := &agent{conn: conn, gate: gate, id: next}
			next = next + 1
			if gate.OnNewAgent != nil {
				gate.OnNewAgent(a)
			}
			return a
		}
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	id 		int64
	conn     Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			if err.Error() != "EOF" {
				log.Debug("read message: %v", err)
			}
			break
		}
		if a.gate.Handle(data, a) == false {
			break
		}
	}
}

func (a *agent) ID() int64{
	return a.id
}

func (a *agent) OnClose() {
	if a.gate.OnCloseAgent != nil {
		a.gate.OnCloseAgent(a)
	}
}

func (a *agent) WriteMsg(data []byte) {
	err := a.conn.WriteMsg(data)
	if err != nil {
		log.Error("write message %v error: %v", data, err)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
