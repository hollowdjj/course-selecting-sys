package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
)

type Gin struct {
	C *gin.Context
}

//发送http响应报文
func (g *Gin) Response(httpCode int, errCode e.ErrNo, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})
}
