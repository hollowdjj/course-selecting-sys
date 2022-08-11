package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
)

type Gin struct {
	C *gin.Context
}

//send http response
func (g *Gin) Response(httpCode int, errCode constval.ErrNo, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  constval.GetErrCodeMsg(errCode),
		"data": data,
	})
}
