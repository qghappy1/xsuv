package etcdv3

import (
	"context"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/qghappy1/xsuv/util/log"
)

type Client = clientv3.Client

type Kv struct {
	K string
	V string
}

func NewKv(k, v string) *Kv {
	return &Kv{k, v}
}

type EtcdC struct {
	Path   string
	mtx    *sync.RWMutex
	nodes  map[string]string
	client *clientv3.Client
}

func NewEtcdV3(addrs []string, path string) *EtcdC {
	this := new(EtcdC)
	this.Path = path
	this.mtx = &sync.RWMutex{}
	this.nodes = make(map[string]string)
	var err error
	this.client, err = clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return this
}

func (this *EtcdC) Close() {
	this.client.Close()
}

func (this *EtcdC) put(resp *clientv3.LeaseGrantResponse, key, val string, timeout int) bool {
	if resp == nil {
		_, err := this.client.Put(context.TODO(), key, val)
		if err != nil {
			log.Error(err.Error())
			return false
		}
		return true
	}
	_, err := this.client.Put(context.TODO(), key, val, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return true
}

func (this *EtcdC) Put(key, val string, timeout int) bool {
	if timeout > 1 {
		resp, err := this.client.Grant(context.TODO(), int64(timeout))
		if err != nil {
			log.Error(err.Error())
			return false
		}
		return this.put(resp, key, val, timeout)
	}
	return this.put(nil, key, val, timeout)
}

func (this *EtcdC) PutLoop(timeout int, kvs ...*Kv) {
	if timeout > 1 {
		resp, err := this.client.Grant(context.TODO(), int64(timeout))
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, kv := range kvs {
			this.put(resp, kv.K, kv.V, timeout)
		}
		go func() {
			ch, err1 := this.client.KeepAlive(context.TODO(), resp.ID)
			if err1 != nil {
				log.Error(err1.Error())
				return
			}
			<-ch
			for {
				for _, kv := range kvs {
					this.put(resp, kv.K, kv.V, timeout)
				}
				time.Sleep(time.Duration(timeout-1) * time.Second)
			}
		}()
	} else {
		log.Fatal("timeout must greater then 1")
	}
}

func (this *EtcdC) Watch() {
	go func() {
		rch := this.client.Watch(context.Background(), this.Path)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					this.mtx.Lock()
					this.nodes[string(ev.Kv.Key)] = string(ev.Kv.Key)
					this.mtx.Unlock()
				case clientv3.EventTypeDelete:
					this.mtx.Lock()
					delete(this.nodes, string(ev.Kv.Key))
					this.mtx.Unlock()
				default:

				}
			}
		}
	}()
}

func (this *EtcdC) watchTest() {
	go func() {
		rch := this.client.Watch(context.Background(), this.Path)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					this.mtx.Lock()
					this.nodes[string(ev.Kv.Key)] = string(ev.Kv.Key)
					this.mtx.Unlock()
					log.Debug("Put %v:%v", string(ev.Kv.Key), string(ev.Kv.Value))
				case clientv3.EventTypeDelete:
					this.mtx.Lock()
					delete(this.nodes, string(ev.Kv.Key))
					this.mtx.Unlock()
					log.Debug("Del %v:%v", string(ev.Kv.Key), string(ev.Kv.Value))
				default:

				}
			}
		}
	}()
}

func (this *EtcdC) Get(key string) map[string]string {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	m := make(map[string]string)
	for k, v := range this.nodes {
		m[k] = v
	}
	return m
}

func ConnectEtcdV3(addrs []string, path string) *Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return client
}

func Get(c *Client, key string, timeout int) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	resp, err := c.Get(ctx, key)
	cancel()
	if err != nil {
		log.ErrorDepth(2, "%v", err)
		return ""
	}
	for _, ev := range resp.Kvs {
		return string(ev.Key)
	}
	return ""
}

func Put(c *Client, key, val string) bool {
	_, err := c.Put(context.TODO(), key, val)
	if err != nil {
		log.ErrorDepth(2, err.Error())
		return false
	}
	return true
}
