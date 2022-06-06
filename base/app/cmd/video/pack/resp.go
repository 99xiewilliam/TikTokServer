package pack

import (
	"errors"
	"go-tiktok/app/kitex_gen/video"
	"go-tiktok/app/pkg/errno"
)

func BuildBaseResp(err error) *video.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *video.BaseResp {
	return &video.BaseResp{StatusCode: err.ErrCode, StatusMsg: err.ErrMsg}
}
