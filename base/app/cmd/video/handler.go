package main

import (
	"context"
	"go-tiktok/app/cmd/video/pack"
	"go-tiktok/app/cmd/video/service"
	"go-tiktok/app/kitex_gen/video"
	"go-tiktok/app/pkg/conf"
	"go-tiktok/app/pkg/errno"
	"time"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	svc *service.Service
}

func NewVideoService(c *conf.Config) video.VideoService {
	s := new(VideoServiceImpl)
	s.svc = service.New(c)
	return s
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	resp = new(video.FeedResponse)

	userId := req.UserId
	latestTime := req.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}
	videoList, nextTime, err := s.svc.FeedVideos(ctx, userId, latestTime)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.VideoList = videoList
	resp.NextTime = nextTime
	return resp, nil
}

// PubAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PubAction(ctx context.Context, req *video.PubActionRequest) (resp *video.PubActionResponse, err error) {
	resp = new(video.PubActionResponse)

	data := req.Data

	err = s.svc.PubVideo(ctx, req.UserId, data, req.Title)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// PubList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PubList(ctx context.Context, req *video.PubListRequest) (resp *video.PubListResponse, err error) {
	resp = new(video.PubListResponse)

	userId := req.UserId
	authorId := req.AuthorId
	videoList, err := s.svc.GetPub(ctx, userId, authorId)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.VideoList = videoList
	return resp, nil
}
