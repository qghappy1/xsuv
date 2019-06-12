package module

import (
	"fmt"
	"time"
	"strings"
	"runtime"
	"xsuv/util/log"
)

type tFunc struct {
	debug	string
	t 		int64
	f		func()
}

type Frame struct {
	*dispatcher
	maxFunSize		int
	frameFun		func()
	frameTime		int
	funcs 			chan *tFunc
	closeSig 		chan bool
}

// frameTime 33ms
func NewFrame(frameTime int, frameFun func()) *Frame {
	f := new(Frame)
	f.frameTime = frameTime
	f.maxFunSize = 1024*10
	f.frameFun = frameFun
	f.funcs = make(chan *tFunc, f.maxFunSize)
	f.dispatcher = newDispatcher(f.maxFunSize)
	f.closeSig = make(chan bool)
	f.run()
	return f
}

func (this *Frame) run(){
	var lastFunc *tFunc
	go func(){
		defer log.ErrorPanic()
		if this.frameTime>0 {
			tick := time.NewTicker(time.Duration(this.frameTime)*time.Millisecond)
			for {
				select {
				case <- this.closeSig:
					return
				case <-tick.C:
					if this.frameFun != nil {
						this.frameFun()
					}
				case cb := <-this.dispatcher.ChanTimer:
					cb.Cb()
				case f, ok := <- this.funcs:
					if ok && f != nil && f.f != nil {
						lastFunc = f
						f.f()
					}
				}
			}
		}else{
			for {
				select {
				case <- this.closeSig:
					return
				case cb := <-this.dispatcher.ChanTimer:
					cb.Cb()
				case f, ok := <- this.funcs:
					if ok && f != nil && f.f != nil {
						lastFunc = f
						f.f()
					}
				}
			}
		}
	}()
}

func (this *Frame) Close() {
	this.closeSig <- true
}

func (this *Frame) Post(f func()){
	tf := new(tFunc)
	tf.f = f

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	tf.debug = fmt.Sprintf("%s.%d", file, line)

	this.funcs <- tf
}

func (this *Frame) PostSync(f func()[]byte)[]byte{
	c := make(chan []byte, 1)
	tf := new(tFunc)
	tf.f = func(){
		c<-f()
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	tf.debug = fmt.Sprintf("%s.%d", file, line)

	this.funcs <- tf
	ret := <-c
	return ret
}
