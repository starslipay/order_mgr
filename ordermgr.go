package main

import (
	"flag"
	"fmt"

	"github.com/starslipay/order_mgr/internal/config"
	"github.com/starslipay/order_mgr/internal/metrics"
	"github.com/starslipay/order_mgr/internal/server"
	"github.com/starslipay/order_mgr/internal/svc"
	"github.com/starslipay/order_mgr/order_mgr_pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ordermgr.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order_mgr_pb.RegisterOrderMgrServer(grpcServer, server.NewOrderMgrServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(metrics.UnaryMetricInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
