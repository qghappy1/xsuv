package module

import (
	"fmt"
)

type _A struct {
	i	int
}

func exampleSignal(){
	OnLogin := "login"
	signal := NewSignal()
	signal.Register(OnLogin, func(args ...interface{}){
		a := args[0].(*_A)
		fmt.Printf("%v\n", a.i)
	})
	a := &_A{1}
	signal.Trigger(OnLogin, a)
}

