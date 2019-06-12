package natslb

import (
	"fmt"
	"time"
	"testing"
	"crypto/md5"
	"xsuv/etcdv3"
	"xsuv/util/log"
)

var serviceType = "game"
//var etcdAddr = "http://111.230.46.154:2379"
var etcdAddr = "http://127.0.0.1:2379"

func natslbCTest(c *etcdv3.Client, key string, i int){
	watcher := NewWatcherKetama(serviceType, c)
	watcher.Start()
	time.Sleep(10*time.Second)
	if s, err := watcher.Get(key); err==nil {
		log.Debug("c:%02d %v:%v:%v", i, time.Now().Unix(), key, s)
	}else{
		log.Debug("%v:%v:error", time.Now().Unix(), key)
	}
}

func natslbSTest(c *etcdv3.Client, serviceName string){
	lb := NewNatsLb(c)
	lb.Register(serviceType, serviceName)
	time.Sleep(30*time.Second)
	lb.UnRegister()
}

func md5_(n string) string {
	//h := md5.New()
	//h.Write([]byte(n))
	return fmt.Sprintf("%x", md5.Sum([]byte(n)))
	//return string(h.Sum(nil))
}

// go test -v xsuv\etcd_mirror
func Test_Etcd(t *testing.T){
	c := ConnectEtcd(etcdAddr)
	for i := 0; i<5; i++ {
		//c := ConnectEtcd(etcdAddr)
		go natslbSTest(c, fmt.Sprintf("game%v", i))
	}
	log.Debug("time:%v", time.Now().Unix())
	for k := 0; k<5; k++ {
		//c := ConnectEtcd(etcdAddr)
		for i := 0; i<5; i++ {
			z := k*5+i
			//log.Debug("1z:%v", z)
			//go natslbCTest(c, md5_(fmt.Sprintf("gate%v", i)), z)
			go natslbCTest(c, fmt.Sprintf("gate%v", i), z)
		}
		time.Sleep(1*time.Second)
	}
	time.Sleep(30*time.Second)
	return	
}


