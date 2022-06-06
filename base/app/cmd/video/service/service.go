package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	uuid "github.com/satori/go.uuid"
	"go-tiktok/app/cmd/video/dao"
	"go-tiktok/app/cmd/video/model/archive"
	"go-tiktok/app/kitex_gen/video"
	"go-tiktok/app/pkg/conf"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Service struct {
	c         *conf.Config
	arc       *dao.Dao
	uploadUrl string
}

func New(c *conf.Config) *Service {
	s := &Service{
		c:         c,
		arc:       dao.New(c),
		uploadUrl: fmt.Sprintf("http://%s/%s/upload", c.Dfs.Host, c.Dfs.Group),
	}
	return s
}

func (s *Service) FeedVideos(ctx context.Context, userId, latestTime int64) (*video.VideoList, int64, error) {
	return s.arc.GetVideos(ctx, userId, latestTime)
}

func (s *Service) PubVideo(ctx context.Context, userId int64, data []byte, title string) error {
	res, err := s.postFile(ctx, data)
	if err != nil {
		return err
	}
	log.Printf("res:%+v\n", res)
	ct := time.Now()
	v := &archive.Video{
		Aid:           userId,
		Title:         title,
		PlayUrl:       res.PlayUrl,
		CoverUrl:      res.CoverUrl,
		PlayMd5:       res.PlayMd5,
		CoverMd5:      res.CoverMd5,
		Status:        1,
		FavoriteCount: 0,
		CommentCount:  0,
		CreateTime:    ct,
		UpdateTime:    ct,
	}
	return s.arc.StoreVideo(ctx, v)
}

type postFileRes struct {
	PlayUrl  string
	PlayMd5  string
	CoverUrl string
	CoverMd5 string
}

type postMsg struct {
	Data struct {
		Domain  string
		Md5     string
		Mtime   int64
		Path    string
		Retcode int
		Retmsg  string
		Scene   string
		Scenes  string
		Size    int
		Src     string
		Url     string
	}
	Message string
	Status  string
}

func (s *Service) postFile(ctx context.Context, data []byte) (*postFileRes, error) {
	var err error

	videoName := fmt.Sprintf("%s.mp4", uuid.NewV4())
	// 当前工作目录是根目录:go-tiktok
	videoPath := "app/temp/" + videoName
	tempVideo, err := os.Create(videoPath)
	if err != nil {
		klog.Errorf("create err:%v\n", err)
		return nil, err
	}
	_, err = tempVideo.Write(data)
	tempVideo.Close()
	if err != nil {
		klog.Errorf("write err:%v\n", err)
		return nil, err
	}

	imageName := fmt.Sprintf("%s.jpg", uuid.NewV4())
	imagePath := "app/temp/" + imageName
	command := exec.Command("ffmpeg", "-i", videoPath, "-y", "-f", "image2", "-ss", "1",
		"-t", "0.001", "-s", "1080x1920", imagePath)
	err = command.Run()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = os.Remove(videoPath)
		if err != nil {
			klog.Debugf("remove temp video fail, name: %s, err: %v\n", videoName, err)
		}
		err = os.Remove(imagePath)
		if err != nil {
			klog.Debugf("remove temp image fail, name: %s, err: %v\n", imageName, err)
		}
	}()

	payload1 := &bytes.Buffer{}
	writer1 := multipart.NewWriter(payload1)
	reader1 := bytes.NewReader(data)
	part1, err := writer1.CreateFormFile("file", videoName)
	_, err = io.Copy(part1, reader1)
	if err != nil {
		return nil, err
	}
	_ = writer1.WriteField("filename", videoName)
	_ = writer1.WriteField("output", "json2")
	writer1.Close()

	payload2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(payload2)
	reader2, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	part2, err := writer2.CreateFormFile("file", imageName)
	_, err = io.Copy(part2, reader2)
	if err != nil {
		return nil, err
	}
	_ = writer2.WriteField("filename", imageName)
	_ = writer2.WriteField("output", "json2")
	reader2.Close()
	writer2.Close()

	client := &http.Client{}
	// 上传视频
	req1, err := http.NewRequest("POST", s.uploadUrl, payload1)
	if err != nil {
		return nil, err
	}
	req1.Header.Set("Content-Type", writer1.FormDataContentType())
	resp1, err := client.Do(req1)
	if err != nil {
		return nil, err
	}
	body1, err := io.ReadAll(resp1.Body)
	resp1.Body.Close()
	if err != nil {
		return nil, err
	}

	var msg1 postMsg
	fmt.Println(string(body1))
	err = json.Unmarshal(body1, &msg1)
	if err != nil {
		return nil, err
	}
	if msg1.Status != "ok" {
		return nil, fmt.Errorf("request dfs fail\n")
	}
	msg1.Data.Url = msg1.Data.Url[:len(msg1.Data.Url)-1] + "0"

	// 上传封面
	req2, err := http.NewRequest("POST", s.uploadUrl, payload2)
	if err != nil {
		return nil, err
	}
	req2.Header.Set("Content-Type", writer2.FormDataContentType())
	resp2, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	body2, err := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	if err != nil {
		return nil, err
	}

	var msg2 postMsg
	err = json.Unmarshal(body2, &msg2)
	if err != nil {
		return nil, err
	}
	if msg2.Status != "ok" {
		return nil, fmt.Errorf("request dfs fail\n")
	}
	msg2.Data.Url = msg2.Data.Url[:len(msg2.Data.Url)-1] + "0"

	return &postFileRes{
		PlayUrl:  msg1.Data.Url,
		PlayMd5:  msg1.Data.Md5,
		CoverUrl: msg2.Data.Url,
		CoverMd5: msg2.Data.Md5,
	}, nil
}

func (s *Service) GetPub(ctx context.Context, userId, authorId int64) (*video.VideoList, error) {
	return s.arc.GetPubById(ctx, userId, authorId)
}
