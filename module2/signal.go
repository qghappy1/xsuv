package module2

import (
	"github.com/qghappy1/xsuv/util/log"
)

type signalT struct {
	sigs map[string][]funName
}

func NewSignal() *signalT {
	sig := new(signalT)
	sig.sigs = make(map[string][]funName)
	return sig
}

type funName struct {
	F    func(...interface{})
	Name string
}

func (this *signalT) Register(sigName string, f func(args ...interface{})) {
	if funcs, ok := this.sigs[sigName]; ok {
		funcs = append(funcs, funName{f, ""})
		this.sigs[sigName] = funcs
	} else {
		funcs := make([]funName, 0)
		funcs = append(funcs, funName{f, ""})
		this.sigs[sigName] = funcs
	}
}

func (this *signalT) RegisterDebug(sigName string, f func(args ...interface{}), fName string) {
	if funcs, ok := this.sigs[sigName]; ok {
		log.Debug("register F:%v", fName)
		funcs = append(funcs, funName{f, fName})
		this.sigs[sigName] = funcs
	} else {
		log.Debug("register F:%v", fName)
		funcs := make([]funName, 0)
		funcs = append(funcs, funName{f, fName})
		this.sigs[sigName] = funcs
	}
}

func (this *signalT) Trigger(sigName string, args ...interface{}) {
	if funcs, ok := this.sigs[sigName]; ok {
		for _, f := range funcs {
			//if sigName == OnLogin {
			//	log.Debug("trigger F:%v", f.Name)
			//}
			f.F(args...)
		}
	}
}
