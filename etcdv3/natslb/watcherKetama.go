package natslb

import (
	"fmt"
	"time"
	"strings"
	"context"
	"xsuv/etcdv3"
	"xsuv/util/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type WatcherKetama struct {
	hash *Ketama
	serviceType string
	client *clientv3.Client
}

func NewWatcherKetama(serviceType string, c *etcdv3.Client) *WatcherKetama {
	w := new(WatcherKetama)
	w.serviceType = serviceType
	w.client = c
	w.hash = NewKetama(10, nil)
	return w
}

func (w *WatcherKetama) Get(key string) (string, error) {
	for i:=0; i<10; i++ {
		if s, ok := w.hash.Get(key); ok {
			return s, nil
		}
		time.Sleep(1*time.Second)
	}
	return "", fmt.Errorf("not find service")
}

func (w *WatcherKetama) watch() {
	service := fmt.Sprintf("/%s/%s", prefix, w.serviceType)
	rch := w.client.Watch(context.Background(), service, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				w.hash.Add(string(ev.Kv.Value))
				//log.Debug("time:%v watch add :%v, key:%v", time.Now().Unix(), string(ev.Kv.Value), service)
			case mvccpb.DELETE:
				v := strings.Replace(string(ev.Kv.Key), service+"/", "", -1)
				w.hash.Remove(v)
				//log.Debug("time:%v watch del :%v, key:%v", time.Now().Unix(), v, string(ev.Kv.Key))
			}
		}
	}
}

func (w *WatcherKetama) Start() {
	go func(){
		defer log.ErrorPanic()
		w.watch()
	}()
}