
package util

import (
	"fmt"
	"testing"
)

// go test -v xsuv\util
func Test_GetAddress(t *testing.T){
	fmt.Println(GetExternalAddress())
	fmt.Println(GetInternalAddress())
}