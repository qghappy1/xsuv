package etcdv3

import (
	"testing"
	"time"
)

var (
	EtcAddress		= []string{"http://127.0.0.1:2379"}
	EtcdPath		= "game"
	ServiceName		= "logdb"
)

// go test -v xsuv\etcd_mirror
func Test_Etcd(t *testing.T){
	etcd := NewEtcdV3(EtcAddress, EtcdPath)
	etcd.watchTest()
	etcd.PutLoop(5, NewKv(EtcdPath+"/game1", "1.0"))
	etcd.Put(EtcdPath+"/version", "1.0", 0)
	time.Sleep(20*time.Second)
	return	
}


