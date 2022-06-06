package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gomodule/redigo/redis"
	"go-tiktok/app/cmd/video/model/archive"
	mredis "go-tiktok/app/cmd/video/model/redis"
	"go-tiktok/app/kitex_gen/video"
	"strconv"
)

func (d *Dao) GetVideos(ctx context.Context, userId, score int64) (*video.VideoList, int64, error) {
	var (
		videoList []*video.Video
		err       error
	)

	conn := d.redis.Get()
	videos, err := redis.Strings(conn.Do("ZREVRANGEBYSCORE", "videos", score, 0, "LIMIT", 0, 20))
	if err != nil {
		return nil, 0, err
	}

	var (
		rv       mredis.Video
		nextTime int64
	)
	for _, v := range videos {
		err = json.Unmarshal([]byte(v), &rv)
		if err != nil {
			return nil, 0, err
		}
		nv := &video.Video{
			Id:            rv.ID,
			PlayUrl:       rv.PlayUrl,
			CoverUrl:      rv.CoverUrl,
			FavoriteCount: rv.FavoriteCount,
			CommentCount:  rv.CommentCount,
			Title:         rv.Title,
			Author: &video.Author{
				Id:            rv.Author.ID,
				Name:          rv.Author.Name,
				FollowCount:   rv.Author.FollowCount,
				FollowerCount: rv.Author.FollowerCount,
			},
		}
		nv.IsFavorite, err = d.IsFavorite(userId, nv.Id)
		if err != nil {
			klog.Errorf("pkg:dao call:IsFavorite err:%v\n", err)
		}
		nv.Author.IsFollow, err = d.IsFollow(userId, nv.Author.Id)
		if err != nil {
			klog.Errorf("pkg:dao call:IsFollow err:%v\n", err)
		}
		videoList = append(videoList, nv)
		nextTime = rv.Score
	}
	return &video.VideoList{Videos: videoList}, nextTime, nil
}

func (d *Dao) StoreVideo(ctx context.Context, video *archive.Video) error {
	var err error
	err = d.db.Table("tb_video").Create(video).Error
	if err != nil {
		return err
	}
	var res struct {
		Id int64 `json:"id"`
	}
	err = d.db.Table("tb_video").
		Select("id").
		Where("aid = ?", video.Aid).
		Scan(&res).Error
	if err != nil {
		return err
	}

	conn := d.redis.Get()
	bytes, err := redis.Bytes(conn.Do("HGET", "author", video.Aid))
	if err != nil {
		return err
	}
	var author mredis.Author
	err = json.Unmarshal(bytes, &author)
	if err != nil {
		return err
	}
	mv := mredis.Video{
		ID:            res.Id,
		Author:        author,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         video.Title,
		Score:         video.CreateTime.Unix(),
	}
	data, err := json.Marshal(mv)
	if err != nil {
		return err
	}
	_, err = redis.Int(conn.Do("ZADD", "videos", mv.Score, string(data)))
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) GetPubById(ctx context.Context, userId, authorId int64) (*video.VideoList, error) {
	var (
		videoList []*video.Video
		err       error
		isFollow  bool
	)

	isFollow, err = d.IsFollow(userId, authorId)
	if err != nil {
		klog.Errorf("pkg:dao call:IsFollow err:%v\n", err)
	}

	rows, err := d.db.Table("tb_video").
		Select("tb_video.id as id, play_url, cover_url, favorite_count, comment_count, title, tb_user.id as aid, username as name, follow_count, follower_count").
		Joins("join tb_user on tb_video.aid=tb_user.id").
		Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var nv video.Video
		var sv struct {
			Id            int64  `json:"id"`
			PlayUrl       string `json:"play_url"`
			CoverUrl      string `json:"cover_url"`
			FavoriteCount int64  `json:"favorite_count"`
			CommentCount  int64  `json:"comment_count"`
			IsFavorite    bool   `json:"is_favorite"`
			Title         string `json:"title"`
			Aid           int64  `json:"aid"`
			Name          string `json:"name"`
			FollowCount   int64  `json:"follow_count"`
			FollowerCount int64  `json:"follower_count"`
			IsFollow      bool   `json:"is_follow"`
		}
		err = d.db.ScanRows(rows, &sv)
		if err != nil {
			klog.Errorf("pkg:dao call:ScanRows err:%v\n", err)
			continue
		}
		nv = video.Video{
			Id:            sv.Id,
			PlayUrl:       sv.PlayUrl,
			CoverUrl:      sv.CoverUrl,
			FavoriteCount: sv.FavoriteCount,
			CommentCount:  sv.CommentCount,
			Title:         sv.Title,
			Author: &video.Author{
				Id:            sv.Aid,
				Name:          sv.Name,
				FollowCount:   sv.FollowCount,
				FollowerCount: sv.FollowerCount,
				IsFollow:      isFollow,
			},
		}
		nv.IsFavorite, err = d.IsFavorite(userId, nv.Id)
		if err != nil {
			klog.Errorf("pkg:dao call:IsFavorite err:%v\n", err)
		}
		videoList = append(videoList, &nv)
	}

	return &video.VideoList{Videos: videoList}, nil
}

func (d *Dao) IsFollow(userId, authorId int64) (bool, error) {
	if userId == 0 {
		return false, nil
	}
	conn := d.redis.Get()
	key := fmt.Sprintf("%d:follow", userId)
	member := strconv.Itoa(int(authorId))
	return redis.Bool(conn.Do("SISMEMBER", key, member))
}

func (d *Dao) IsFavorite(userId, id int64) (bool, error) {
	if userId == 0 {
		return false, nil
	}
	conn := d.redis.Get()
	key := fmt.Sprintf("%d:favorite", userId)
	member := strconv.Itoa(int(id))
	return redis.Bool(conn.Do("SISMEMBER", key, member))
}
