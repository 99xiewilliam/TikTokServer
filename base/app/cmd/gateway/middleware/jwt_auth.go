package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-tiktok/app/pkg/errno"
	"net/http"
	"time"
)

var (
	jwtKey = []byte("jwt-key-f02ea730-d702-423a-9cbe-39aec9c1e257")
)

type AuthClaims struct {
	UserId int64
	jwt.StandardClaims
}

func GenToken(userId int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)
	claims := AuthClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func parseToken(token string) (*AuthClaims, bool) {
	tokenObj, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, _ := tokenObj.Claims.(*AuthClaims); tokenObj.Valid {
		return claims, true
	} else {
		return nil, false
	}
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.JSON(http.StatusOK, gin.H{
				"status_code": errno.NeedLoginErrCode,
				"status_msg":  "用户未登录, 请先登录!",
			})
			c.Abort()
			return
		}
		token, ok := parseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"status_code": errno.TokenValidErrCode,
				"status_msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > token.ExpiresAt {
			c.JSON(http.StatusOK, gin.H{
				"status_code": errno.TokenExpiredErrCode,
				"status_msg":  "Token过期",
			})
			c.Abort()
			return
		}
		c.Set("user_id", token.UserId)
		c.Next()
	}
}

func GetUserId(c *gin.Context) int64 {
	token := c.Query("token")
	if claims, ok := parseToken(token); ok {
		return claims.UserId
	} else {
		return 0
	}
}
