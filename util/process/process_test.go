package process

import (
	"time"
	"testing"
)

func Test_Process(t *testing.T){
	p := NewProcess()
	ret := make(chan bool)
	p.Post(func(){
		t.Log("hello world")
		ret <- true
		panic(0)
	})

	t.Logf("ret:%v", <-ret)
	p.Close()
	p.Post(func(){
		t.Log("hello golang")
	})
	time.Sleep(time.Second*1)
}

