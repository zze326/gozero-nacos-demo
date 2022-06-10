package main

import (
	"flag"
	"fmt"

	commonConfig "demo/common/config"
	"demo/microservice/user/rpc/internal/config"
	"demo/microservice/user/rpc/internal/server"
	"demo/microservice/user/rpc/internal/svc"
	"demo/microservice/user/rpc/types/app"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	_ "github.com/zze326/zero-contrib/zrpc/registry/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "microservice/user/rpc/etc/nacos.yaml", "the config file")

func main() {
	flag.Parse()

	var c = new(config.Config)
	commonConfig.MustRegister(commonConfig.MustLoad(*configFile, c), &c.RpcServerConf)
	ctx := svc.NewServiceContext(c)
	svr := server.NewUserServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		app.RegisterUserServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
