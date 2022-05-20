package main

import (
	"github.com/TikTokServer/rpc"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	rpc.InitRPC()
	initRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
