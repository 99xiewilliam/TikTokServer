package router

import (
	"github.com/gin-gonic/gin"
	"go-tiktok/app/cmd/gateway/controller"
	"go-tiktok/app/cmd/gateway/middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	g := router.Group("/douyin")
	{
		//g.POST("/user/login/", controller.Login)
		//g.GET("/user/", controller.UserInfo)
		g.GET("/feed", controller.Feed)
		g.Use(middleware.JwtAuth())
		g.POST("/publish/action/", controller.PubAction)
		g.GET("/publish/list/", controller.PubList)
		//g.GET("/favorite/list/", controller.FavoriteList)
	}

	return router
}
