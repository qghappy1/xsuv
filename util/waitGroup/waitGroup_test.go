
package waitGroup

import (
    "fmt"
	"testing"
)
func test(){
	GoWrap(func(){
		i := 0
		z := 0
		i = 1/z+1
		i = i+1
	})
}

func Test_WaitGroup(t *testing.T){
	test()

	if IsSigStop() {
		fmt.Printf("TestExit.exit \n")
		return 
	}
}