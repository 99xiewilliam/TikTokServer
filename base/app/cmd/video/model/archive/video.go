package archive

import "time"

type Video struct {
	ID            int64     `json:"id"`
	Aid           int64     `json:"aid"`
	Title         string    `json:"title"`
	PlayUrl       string    `json:"play_url"`
	CoverUrl      string    `json:"cover_url"`
	PlayMd5       string    `json:"play_md5"`
	CoverMd5      string    `json:"cover_md5"`
	Status        int16     `json:"status"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}
