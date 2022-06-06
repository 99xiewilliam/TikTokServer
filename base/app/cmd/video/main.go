package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	video "go-tiktok/app/kitex_gen/video/videoservice"
	"go-tiktok/app/pkg/bound"
	"go-tiktok/app/pkg/conf"
	"go-tiktok/app/pkg/constants"
	"go-tiktok/app/pkg/middleware"
	"net"
)

func main() {
	c := conf.Default()

	r, err := etcd.NewEtcdRegistry([]string{c.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	//address := os.Args[1]
	//fmt.Println(address)
	//if address == "" {
	//	panic("address is none")
	//}

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	svr := video.NewServer(NewVideoService(c),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.VideoServiceName}),
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithSuite(trace.NewDefaultServerSuite()),
		server.WithBoundHandler(bound.NewCpuLimitHandler()),
		server.WithRegistry(r),
	)

	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
