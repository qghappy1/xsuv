
package nats

import (
	"log"
	"time"
	"testing"
)

func handle(msg []byte) []byte {
	log.Printf("arg1:%v", string(msg))
	time.Sleep(2*time.Second)
	return []byte("z1")
}

func handle2(msg []byte) []byte {
	log.Printf("arg2:%v", string(msg))
	time.Sleep(2*time.Second)
	return []byte("z2")
}

func Publish(nc, nc2 *Nats){
	subj := "person"
	nc.Subscribe(subj, handle)
	nc2.Subscribe(subj, handle2)

	nc.Publish(subj, []byte("sb"))
	time.Sleep(1*time.Second)
}

func RPC(nc *Nats){
	subj := "person2"
	nc.Register(subj, handle)

	ret, err2 := nc.Call(subj, []byte("sb"), 5000)
	if err2 != nil {
		log.Fatal("call err:", err2)
	}
	log.Printf("Recv RPC: '%v'\n", string(ret))
	time.Sleep(3*time.Second)
	ret, err2 = nc.Call(subj, []byte("sb222"), 5000)
	if err2 != nil {
		log.Fatal("call err:", err2)
	}
	log.Printf("Recv RPC: '%v'\n", string(ret))	
		
	time.Sleep(1*time.Second)
}


func Test_Publish(t *testing.T){
	urls := "nats://0.0.0.0:4242"
	username := "derek"
	password := "T0pS3cr3t"
	nc := NewNats(urls, username, password)
	nc2 := NewNats(urls, username, password)
	RPC(nc)
	Publish(nc, nc2)
}
