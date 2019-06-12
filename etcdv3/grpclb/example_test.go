package grpclb

import (
	"fmt"
	"time"
	"testing"
)



// go test -v xsuv\etcd_mirror
func Test_GRPC(t *testing.T){
	etcd := []string{"http://111.230.46.154:2379"}
	ports := []int{9000, 9001, 9002}
	fmt.Println("test")
	for i := 0; i<3; i++ {
		go ExampleStartService(etcd, fmt.Sprint(i), ports[i])
	}
	go ExampleStreamClient2(etcd)
	go ExampleStreamClient2(etcd)
	time.Sleep(20*time.Second)
	return	
}


