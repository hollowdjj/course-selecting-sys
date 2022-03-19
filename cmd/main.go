package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/gredis"
	"github.com/hollowdjj/course-selecting-sys/pkg/setting"
	"github.com/hollowdjj/course-selecting-sys/routers"
	"log"
	"net/http"
)

func init() {
	setting.Setup()
	models.Setup()
	gredis.SetUp()
}

func main() {
	gin.SetMode(setting.AppSetting.RunMode)
	router := routers.RegisterRouter()

	addr := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	server := &http.Server{
		Addr:           addr,         //监听的端口号，格式为:8000
		Handler:        router,       //http句柄，用于处理程序响应HTTP请求
		ReadTimeout:    readTimeout,  //读取超时时间
		WriteTimeout:   writeTimeout, //写超时时间
		MaxHeaderBytes: 1 << 20,      //http报文head的最大字节数
	}

	log.Printf("[info] start http server listening %s", addr)
	server.ListenAndServe()
}
