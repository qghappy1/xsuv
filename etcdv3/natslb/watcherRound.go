package natslb

import (
	"fmt"
	"sync"
	"strings"
	"context"
	"xsuv/etcdv3"
	"xsuv/util/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type IWatcher interface {
	Get(key string) (string, error)
	Start()
}

type WatcherRound struct {
	mtx sync.RWMutex
	nodes []string
	next uint32
	serviceType string
	client *clientv3.Client
}

func NewWatcherRound(serviceType string, c *etcdv3.Client) *WatcherRound {
	w := new(WatcherRound)
	w.serviceType = serviceType
	w.client = c
	w.nodes = make([]string, 0)
	w.mtx = sync.RWMutex{}
	w.next = 0
	return w
}

func (w *WatcherRound) Get(_ string) (string, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	size := uint32(len(w.nodes))
	if size == 0 {
		return "", fmt.Errorf("not exist service")
	}
	w.next = w.next + 1
	idx := int(w.next%size)
	return w.nodes[idx], nil
}

func (w *WatcherRound) watch() {
	service := fmt.Sprintf("/%s/%s", prefix, w.serviceType)
	rch := w.client.Watch(context.Background(), service, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				w.mtx.Lock()
				key := string(ev.Kv.Value)
				exist := false
				for _, v := range w.nodes {
					if v == key {
						exist = true
						break
					}
				}
				if !exist { w.nodes = append(w.nodes, key) }
				w.mtx.Unlock()
			case mvccpb.DELETE:
				key := strings.Replace(string(ev.Kv.Key), service+"/", "", -1)
				w.mtx.Lock()
				for i, v := range w.nodes {
					if v == key {
						w.nodes = append(w.nodes[:i], w.nodes[i+1:]...)
						break
					}
				}
				w.mtx.Unlock()
			}
		}
	}
}

func (w *WatcherRound) Start() {
	go func(){
		defer log.ErrorPanic()
		w.watch()
	}()
}