package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io/ioutil"
	"log"
	"os"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	//获取临时路径
	tempLogDir, err := ioutil.TempDir("", "log")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempLogDir)

	logDir := fmt.Sprintf(
		"%s%snacos%slog",
		tempLogDir,
		string(os.PathSeparator),
		string(os.PathSeparator),
	)
	//临时cache路径
	tempCacheDir, err := ioutil.TempDir("", "log")
	defer os.RemoveAll(tempCacheDir)
	cacheDir := fmt.Sprintf("%s%snacos%scache",
		tempCacheDir,
		string(os.PathSeparator),
		string(os.PathSeparator),
	)
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              logDir,
		CacheDir:            cacheDir,
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Panic(err)
	}

	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	r := nacos.New(client)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(r),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
