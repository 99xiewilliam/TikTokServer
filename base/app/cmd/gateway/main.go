package main

import (
	"go-tiktok/app/cmd/gateway/router"
	"go-tiktok/app/cmd/gateway/rpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rpc.Init()
	router.HttpServerRun()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	router.HttpServerStop()
}
