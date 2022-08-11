package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/bootstrap"
	"github.com/hollowdjj/course-selecting-sys/conf"
	"github.com/hollowdjj/course-selecting-sys/routers"
)

func main() {
	bootstrap.Run()

	gin.SetMode(conf.GetApp().RunMode)
	router := routers.RegisterRouter()

	serverConf := conf.GetServer()
	addr := fmt.Sprintf(":%d", serverConf.HttpPort)
	readTimeout := time.Second * time.Duration(serverConf.ReadTimeout)
	writeTimeout := time.Second * time.Duration(serverConf.WriteTimeout)
	server := &http.Server{
		Addr:           addr,         //监听的端口号，格式为:8000
		Handler:        router,       //http句柄，用于处理程序响应HTTP请求
		ReadTimeout:    readTimeout,  //读取超时时间
		WriteTimeout:   writeTimeout, //写超时时间
		MaxHeaderBytes: 1 << 20,      //http报文head的最大字节数
	}

	log.Printf("start http server listening %s", addr)
	server.ListenAndServe()
}
