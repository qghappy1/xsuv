
package nats

import (
	"time"
	"strings"
	"xsuv/util/log"
	"github.com/nats-io/go-nats"	
)

type Nats struct {
	conn	*nats.Conn
}

func NewNats(urls, username, password string) (*Nats) {
	var err error
	p := new(Nats)
	p.conn, err = connectNats(urls, username, password)
	if err != nil {
		log.Fatal("nats conn fail.err:%v", err)
		return nil
	}
	p.conn.Reconnects = 5
	return p
}

func (this *Nats) Close() {
	this.conn.Close()	
}

// 发布
func (this *Nats) Publish(subj string, msg []byte) error {
	return this.conn.Publish(subj, msg)
}

// rpc call
func (this *Nats) Call(subj string, arg []byte, timeout int) ([]byte, error) {
	msg, err := this.conn.Request(subj, arg, time.Duration(timeout)*time.Millisecond)
	//msg, err := this.enc.Request(subj, msg, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return nil, err
	}
	return msg.Data, err
}

// 
func connectNats(urls, username, password string) (*nats.Conn, error){
	opts := nats.DefaultOptions
	opts.Servers = strings.Split(urls, ",")
	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}
	opts.Secure = false
	opts.User = username
	opts.Password = password
	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	return nc, nil
}

// 订阅
func (this *Nats) Subscribe(subj string, handle func([]byte)[]byte) error {
	_, err := this.conn.Subscribe(subj, func(msg *nats.Msg){
		handle(msg.Data)
	})
	return err
}

// rpc
func (this *Nats) Register(subj string, handle func([]byte)[]byte) error {
	_, err := this.conn.Subscribe(subj, func(msg *nats.Msg){
		ret := handle(msg.Data)	
		this.conn.Publish(msg.Reply, ret)
	})
	return err
}









