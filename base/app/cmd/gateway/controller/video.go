package controller

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"go-tiktok/app/cmd/gateway/middleware"
	"go-tiktok/app/cmd/gateway/rpc"
	"go-tiktok/app/kitex_gen/video"
	"go-tiktok/app/pkg/errno"
	"io"
	"net/http"
	"strconv"
)

func Feed(c *gin.Context) {
	var (
		userId     int64
		latestTime int64
		err        error
	)

	userId = middleware.GetUserId(c)

	latestTimeStr := c.Query("latest_time")
	latestTime, err = strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		latestTime = 0
	}

	ctx := context.Background()

	req := &video.FeedRequest{
		UserId:     userId,
		LatestTime: latestTime,
	}

	resp, err := rpc.Feed(ctx, req)
	if err != nil {
		klog.Errorf("pkg:controller rpc.Feed err:%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"status_code": errno.ServiceErrCode,
			"status_msg":  "service error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": errno.SuccessCode,
		"status_msg":  "feed success",
		"next_time":   resp.NextTime,
		"video_list":  resp.VideoList.Videos,
	})
}

func PubAction(c *gin.Context) {
	var (
		userId int64
		title  string
		data   []byte
		err    error
	)

	if value, exists := c.Get("user_id"); exists {
		userId = value.(int64)
	} else {
		userId = 0
	}

	title = c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": errno.ParamErrCode,
			"status_msg":  "标题不能为空",
		})
		return
	}

	header, _ := c.FormFile("data")
	file, _ := header.Open()
	defer file.Close()
	data, err = io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": errno.ParamErrCode,
			"status_msg":  "视频读取失败",
		})
		return
	}

	ctx := context.Background()

	req := &video.PubActionRequest{
		UserId: userId,
		Data:   data,
		Title:  title,
	}

	err = rpc.PubAction(ctx, req)
	if err != nil {
		klog.Errorf("pkg:controller rpc.PubAction err:%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"status_code": errno.ServiceErrCode,
			"status_msg":  "service error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": errno.SuccessCode,
		"status_msg":  "publish success",
	})
}

func PubList(c *gin.Context) {
	var (
		userId   int64
		authorId int64
		err      error
	)

	value, _ := c.Get("user_id")
	userId = value.(int64)

	str := c.Query("user_id")
	if str == "" || str == "0" {
		authorId = userId
	} else {
		authorId, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": errno.ParamErrCode,
				"status_msg":  "user_id格式错误",
			})
			return
		}
	}

	ctx := context.Background()

	req := &video.PubListRequest{
		UserId:   userId,
		AuthorId: authorId,
	}

	resp, err := rpc.PubList(ctx, req)
	if err != nil {
		klog.Errorf("pkg:controller rpc.PubList err:%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"status_code": errno.ServiceErrCode,
			"status_msg":  "service error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": errno.SuccessCode,
		"status_msg":  "success",
		"video_list":  resp.VideoList.Videos,
	})
}

//func FavoriteList(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"status_code": errno.SuccessCode,
//		"status_msg":  "success",
//		"video_list":  []video.Video{},
//	})
//}
