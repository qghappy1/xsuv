package grpclb

import (
	"os"
	"fmt"
	"log"
	"net"
	"time"
	"syscall"
	"os/signal"
	"xsuv/etcdv3/grpclb/api"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
	registry "github.com/liyue201/grpc-lb/registry/etcd"
	"io"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	Node string
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	fmt.Printf("Node:%v Receive is %s\n", s.Node, in.Name)
	return &api.HelloReply{Message: "Hello " + in.Name}, nil
}

// 双向流
func (s *server) SayStream(stream api.Greeter_SayStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("s read done")
			return nil
		}
		if err != nil {
			fmt.Println("ERR",err)
			return err
		}
		fmt.Println("node ", s.Node, " recv message ",in.Name)
		if err := stream.Send(&api.HelloReply{Message: "Hello " + in.Name + " from "+s.Node}); err != nil {
			return err
		}
	}
	return nil
}

func ExampleStartService(etcdip []string, nodeID string, port int) {
	etcdConfg := etcd.Config{
		Endpoints: etcdip,
	}

	registry, err := registry.NewRegistry(
		registry.Option{
			EtcdConfig:  etcdConfg,
			RegistryDir: "/service",
			ServiceName: "test",
			NodeID:      nodeID,
			NData: registry.NodeData{
				Addr: fmt.Sprintf("0.0.0.0:%d", port),
				//Metadata: map[string]string{"weight": "1"}, //这里配置权重，不配置默认是1
			},
			Ttl: 10 * time.Second,
		})
	if err != nil {
		fmt.Println(1)
		log.Panic(err)
		return
	}
	go registry.Register()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		fmt.Println(3)
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	fmt.Printf("starting test service at %v\n", nodeID)
	s := grpc.NewServer()
	api.RegisterGreeterServer(s, &server{Node:nodeID})
	go s.Serve(lis)

	sig := <-ch
	fmt.Printf("receive signal '%v'\n", sig)
	s.GracefulStop()
	registry.Deregister()
}