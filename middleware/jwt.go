package middleware

import (
	"context"
	"github.com/TikTokServer/controller"
	"github.com/TikTokServer/kitex_gen/user"
	"github.com/TikTokServer/rpc"
	"github.com/gin-gonic/gin"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
)

var identityKey = "id"
var SecretKey = "secret"

func authenticator(c *gin.Context) (interface{}, error) {
	username := c.Query("username")
	password := c.Query("password")
	if len(username) == 0 || len(password) == 0 {
		return "", jwt.ErrMissingLoginValues
	}

	userId, err := rpc.CheckUser(context.Background(), &user.CheckUserRequest{UserName: username, Password: password})
	if err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	c.Set("UserId", userId)
	return userId, nil
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(int64); ok {
		return jwt.MapClaims{
			identityKey: v,
		}
	}
	return jwt.MapClaims{}
}

//func authrizator(data interface{}, c *gin.Context) bool {
//	return true
//}

func JwtMiddlewareInit() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:           []byte(SecretKey),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		PayloadFunc:   payloadFunc,
		Authenticator: authenticator,
		//Authorizator:  authrizator,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		LoginResponse: controller.LoginResponse,
	})
	if err != nil {
		panic(err)
	}
	return authMiddleware, nil
}
