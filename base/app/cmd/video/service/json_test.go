package service

import (
	"fmt"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	//	jsonData := `{
	//  "data": {
	//    "domain": "http://127.0.0.1:8086",
	//    "md5": "ab22aa32da861e828133482652403e33",
	//    "mtime": 1654280209,
	//    "path": "/group1/default/20220604/02/16/5/53bad4bc-d053-4a00-82d1-23e697f22f0f",
	//    "retcode": 0,
	//    "retmsg": "",
	//    "scene": "default",
	//    "scenes": "default",
	//    "size": 1164750,
	//    "src": "/group1/default/20220604/02/16/5/53bad4bc-d053-4a00-82d1-23e697f22f0f",
	//    "url": "http://127.0.0.1:8086/group1/default/20220604/02/16/5/53bad4bc-d053-4a00-82d1-23e697f22f0f?name=53bad4bc-d053-4a00-82d1-23e697f22f0f&download=1"
	//  },
	//  "message": "",
	//  "status": "ok"
	//}`
	//	var msg postMsg
	//	err := json.Unmarshal([]byte(jsonData), &msg)
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Printf("%+v\n", msg)
	//	}

	//author := redis.Author{
	//	ID:            1,
	//	Name:          "zhi",
	//	FollowCount:   0,
	//	FollowerCount: 0,
	//}
	//s, _ := json.Marshal(author)
	//fmt.Println(string(s))

	//var author redis.Author
	//str := "{\"id\":1,\"name\":\"zhi\",\"follow_count\":0,\"follower_count\":0}"
	//err := json.Unmarshal([]byte(str), &author)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Printf("%+v\n", author)
	//}

	fmt.Println(time.Now().Unix())
}
