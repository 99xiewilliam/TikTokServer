package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"go-tiktok/app/kitex_gen/video"
	"go-tiktok/app/kitex_gen/video/videoservice"
	"go-tiktok/app/pkg/conf"
	"go-tiktok/app/pkg/constants"
	"go-tiktok/app/pkg/errno"
	"go-tiktok/app/pkg/middleware"
	"time"
)

var videoClient videoservice.Client

func Init() {
	c := conf.Default()

	r, err := etcd.NewEtcdResolver([]string{c.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	cli, err := videoservice.NewClient(
		constants.VideoServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),
		client.WithRPCTimeout(3*time.Second),
		client.WithConnectTimeout(50*time.Millisecond),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithSuite(trace.NewDefaultClientSuite()),
		client.WithResolver(r))
	if err != nil {
		panic(err)
	}

	videoClient = cli
}

func Feed(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
	resp, err := videoClient.Feed(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.New(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp, nil
}

func PubAction(ctx context.Context, req *video.PubActionRequest) error {
	resp, err := videoClient.PubAction(ctx, req)
	if err != nil {
		return err
	}
	if resp.BaseResp.StatusCode != 0 {
		return errno.New(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return nil
}

func PubList(ctx context.Context, req *video.PubListRequest) (*video.PubListResponse, error) {
	resp, err := videoClient.PubList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != 0 {
		return nil, errno.New(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp, nil
}
