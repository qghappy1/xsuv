package grpclb

import (
	"fmt"
	"log"
	"time"
	etcd "github.com/coreos/etcd/client"
	grpclb "github.com/liyue201/grpc-lb"
	"xsuv/etcdv3/grpclb/api"
	registry "github.com/liyue201/grpc-lb/registry/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"io"
)

// 随机策略
func exampleRandomClient(etcdip []string){
	etcdConfg := etcd.Config{
		Endpoints: etcdip,
	}
	r := registry.NewResolver("/service", "test", etcdConfg)
	b := grpclb.NewBalancer(r, grpclb.NewRandomSelector())
	c, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()


	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := api.NewGreeterClient(c)
		resp, err := client.SayHello(context.Background(), &api.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("Send is %s\n", resp.Message)
		}
	}
}

// 轮循
func exampleRoundClient(etcdip []string){
	etcdConfg := etcd.Config{
		Endpoints: etcdip,
	}
	r := registry.NewResolver("/service", "test", etcdConfg)
	b := grpclb.NewBalancer(r, grpclb.NewRoundRobinSelector())
	c, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := api.NewGreeterClient(c)
		resp, err := client.SayHello(context.Background(), &api.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("Send is %s\n", resp.Message)
		}
	}
}

// 一致性哈希
func ExampleKetamaClient(etcdip []string){
	etcdConfg := etcd.Config{
		Endpoints: etcdip,
	}
	r := registry.NewResolver("/service", "test", etcdConfg)
	b := grpclb.NewBalancer(r, grpclb.NewKetamaSelector(grpclb.DefaultKetamaKey))
	c, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithTimeout(time.Second))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()
	key := 0
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		key = key+1
		hashData := fmt.Sprintf("%d", key)
		client := api.NewGreeterClient(c)
		resp, err := client.SayHello(context.WithValue(context.Background(), grpclb.DefaultKetamaKey, hashData), &api.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("Send is %s\n", resp.Message)
		}
	}
}

// 双向流模式
func ExampleStreamClient(address string){
	c, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()

	i := 1
	waitc := make(chan struct{})
	client := api.NewGreeterClient(c)
	stream, err := client.SayStream(context.Background())
	if err != nil {
		fmt.Println("%v.RouteChat(_) = _, %v", client, err)
	}
	if err := stream.Send(&api.HelloRequest{Name:fmt.Sprintf("world %v", i)}); err != nil {
		fmt.Println("Failed to send a note: %v", err)
	}
	fmt.Println("send message")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("c read done ")
			close(waitc)
			return
		}
		if err != nil {
			fmt.Println("Failed to receive a note", err)
			return
		}
		i = i+1
		fmt.Println("client recv message ",in.Message)
		time.Sleep(time.Second)
		if err := stream.Send(&api.HelloRequest{Name:fmt.Sprintf("world %v", i)}); err != nil {
			fmt.Println("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

// 双向流模式
func ExampleStreamClient2(etcdip []string){
	etcdConfg := etcd.Config{
		Endpoints: etcdip,
	}
	r := registry.NewResolver("/service", "test", etcdConfg)
	b := grpclb.NewBalancer(r, grpclb.NewRandomSelector())
	c, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()

	i := 1
	waitc := make(chan struct{})
	client := api.NewGreeterClient(c)
	stream, err := client.SayStream(context.Background())
	if err != nil {
		fmt.Println("%v.RouteChat(_) = _, %v", client, err)
	}
	if err := stream.Send(&api.HelloRequest{Name:fmt.Sprintf("world %v", i)}); err != nil {
		fmt.Println("Failed to send a note: %v", err)
	}
	fmt.Println("send message")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("c read done ")
			close(waitc)
			return
		}
		if err != nil {
			fmt.Println("Failed to receive a note ", err)
			return
		}
		i = i+1
		fmt.Println("client recv message ",in.Message)
		time.Sleep(time.Second)
		if err := stream.Send(&api.HelloRequest{Name:fmt.Sprintf("world %v", i)}); err != nil {
			fmt.Println("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}




