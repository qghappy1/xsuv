package natslb

import (
	"fmt"
	"time"
	"strings"
	"context"
	"xsuv/nats"
	"xsuv/util/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

type Client = clientv3.Client
type Nats = nats.Nats

var (
	NewNats = nats.NewNats
	prefix = "service"
)

type NatsLb struct {
	client *Client
	serviceKey string
	stopSignal chan bool
}

func NewNatsLb(c *Client) *NatsLb {
	g := new(NatsLb)
	g.client = c
	g.stopSignal = make(chan bool, 1)
	return g
}

func ConnectEtcd(etcdAddrs string) *Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddrs, ","),
	})
	if err != nil {
		log.Fatal("create etcd3 client failed: %v", err)
	}
	return client
}

// Register serviceName:game timeout:10 return game11
func (this *NatsLb) Register(serviceType, serviceName string) {
	this.serviceKey = fmt.Sprintf("/%s/%s/%s", prefix, serviceType, serviceName)
	go func() {
		defer log.ErrorPanic()
		tick := time.NewTicker(6 * time.Second)
		for {
			resp, _ := this.client.Grant(context.TODO(), int64(10))
			_, err := this.client.Get(context.Background(), this.serviceKey)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					if _, err := this.client.Put(context.TODO(), this.serviceKey, serviceName, clientv3.WithLease(resp.ID)); err != nil {
						log.Error("natslb: set service '%s' with ttl to etcd3 failed: %s", serviceName, err.Error())
					}
					//log.Debug("register key:%v, value:%v", this.serviceKey, serviceName)
				} else {
					log.Error("natslb: service '%s' connect to etcd3 failed: %s", serviceName, err.Error())
				}
			} else {
				if _, err := this.client.Put(context.Background(), this.serviceKey, serviceName, clientv3.WithLease(resp.ID)); err != nil {
					log.Error("natslb: refresh service '%s' with ttl to etcd3 failed: %s", serviceName, err.Error())
				}
				//log.Debug("time:%v register key:%v, value:%v", time.Now().Unix(), this.serviceKey, serviceName)
			}
			select {
			case <-this.stopSignal:
				tick.Stop()
				return
			case <- tick.C:
			}
			//log.Error("time:%v register:%v", time.Now().Unix(), serviceName)
		}
	}()
}

// UnRegister delete registered service from etcd
func (this *NatsLb) UnRegister() error {
	this.stopSignal <- true
	var err error;
	if _, err := this.client.Delete(context.Background(), this.serviceKey); err != nil {
		log.Error("natslb: deregister '%s' failed: %s", this.serviceKey, err.Error())
	}
	return err
}