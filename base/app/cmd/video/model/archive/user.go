package archive

import "time"

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	FollowCount   int64     `json:"follow_count"`
	FollowerCount int64     `json:"follower_count"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}
