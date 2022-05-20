package controller

import (
	"context"
	"fmt"
	"github.com/TikTokServer/kitex_gen/user"
	"github.com/TikTokServer/rpc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserParam struct {
	username string
	password string
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Invalid parameter"})
		return
	}

	presence, err := rpc.CheckUserPresence(context.Background(), &user.CheckUserPresenceRequest{UserName: username})
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Internal error"})
		return
	}
	if presence {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User already exists"})
		return
	}
	userId, err := rpc.CreateUser(context.Background(), &user.CreateUserRequest{
		UserName: username,
		Password: password,
	})

	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Internal error"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))
	c.JSON(http.StatusOK, UserLoginResponse{Response: Response{StatusCode: 0}, UserId: userId, Token: tokenString})
}

//func Login(c *gin.Context) {
//	username := c.Query("username")
//	password := c.Query("password")
//
//	token := username + password
//
//	if user, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 0},
//			UserId:   user.Id,
//			Token:    token,
//		})
//	} else {
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
//		})
//	}
//}

func LoginResponse(c *gin.Context, code int, message string, time time.Time) {
	userId, exists := c.Get("UserId")
	if !exists {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userId.(int64),
		Token:    message,
	})
}

func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Invalid parameter"})
		return
	}
	u, err := rpc.GetUserInfo(context.Background(), &user.GetUserInfoRequest{UserId: int64(userIdInt)})
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Internal error"})
	}
	fmt.Printf("%+v\n", u)
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: User{
			Id: u.UserId,
			//Name: u.UserName,
		},
	})
}
