
package exception

import (
	"time"
	"testing"
)

func Test_catch(t *testing.T){
	go func(){
		defer Catch(nil)
		panic(55)
	}()	
	time.Sleep(time.Second*1)	
}

func Test_CatchAndExit(t *testing.T){
	go func(){
		defer CatchAndExit(nil)
		panic(55)
	}()	
	time.Sleep(time.Second*1)	
}

