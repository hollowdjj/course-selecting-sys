package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/conf"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/hollowdjj/course-selecting-sys/proxy"
)

var Router *gin.Engine

func Run() {
	conf.LoadConfig("./conf/conf.ini")

	logger.InitLogger("./log")

	models.InitDb()

	proxy.InitHttpPool()
}
