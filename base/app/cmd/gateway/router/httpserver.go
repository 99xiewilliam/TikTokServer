package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	addr           string
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	addr = os.Args[1]
	if addr == "" {
		panic("addr is none")
	}
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		log.Printf("[INFO] HttpServerRun:%s\n", addr)
		//if err := HttpSrvHandler.ListenAndServeTLS("/home/zhi/crt/server.crt", "/home/zhi/crt/server.key"); err != nil {
		//	log.Fatalf("[ERROR] HttpServerRun:192.168.43.6:8080 err:%v\n", err)
		//}
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf("[ERROR] HttpServerRun:%s err:%v\n", addr, err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf("[INFO] HttpServerStop stopped\n")
}
