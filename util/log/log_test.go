
package log

import (
	"fmt"
	"testing"
)

// go test -v xsuv\util\log
func Test_log(t *testing.T){
	fmt.Println("test log")
	SetFlag(DEBUG)
	v := 1
	Debug("this is debug:%v", v)
	Info("this is info:%v", v)
	Warn("this is info:%v", v)
	Error("this is info:%v", v)
}