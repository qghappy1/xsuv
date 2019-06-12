package process

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/qghappy1/xsuv/util/log"
	"github.com/qghappy1/xsuv/util/waitGroup"
)

type tFunc struct {
	debug string
	f     func()
}

type Process struct {
	funcs chan *tFunc
	isrun bool
	debug string
}

func NewProcess() *Process {
	p := new(Process)
	p.funcs = make(chan *tFunc, 1024*1024)
	p.isrun = false

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	i := strings.LastIndex(file, "/")
	file = file[i+1:]
	p.debug = fmt.Sprintf("%s.%d", file, line)

	p.Run()
	return p
}

func (this *Process) Run() {
	if this.isrun {
		return
	}
	this.isrun = true

	var lastFunc *tFunc
	invoke := func() {
		for {
			if waitGroup.IsSigStop() {
				break
			}
			select {
			case f, ok := <-this.funcs:
				if ok && f != nil && f.f != nil {
					lastFunc = f
					f.f()
				} else {
					return
				}

			default:
				time.Sleep(time.Millisecond * 5)
			}
		}
	}
	finish := func() {
		if lastFunc != nil {
			log.Error("invoke:%v error. process:%v", lastFunc.debug, this.debug)
		}
		this.isrun = false
		this.Run()
	}
	waitGroup.GoWrapEx(invoke, finish)
}

func (this *Process) Post(f func()) {
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

func (this *Process) Close() {
	this.funcs <- nil
}
