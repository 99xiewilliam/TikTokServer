package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestMySql(t *testing.T) {
	_, err := gorm.Open("mysql", "tiktok:123456@tcp(127.0.0.1:3306)/db_tiktok?charset=utf8&parseTime=true&loc=Asia%2FChongqing")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}
}
