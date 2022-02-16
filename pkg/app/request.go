package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"net/http"
)

//MakeErrors 向日志中写入表单验证时发生的错误
func MakeErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}

//根据cookie值判断当前登录的用户是不是admin
func IsAdmin(c *gin.Context) bool {
	cookie, err := c.Cookie("camp-seesion")
	if err == http.ErrNoCookie {
		return false
	}
	userType := utility.GetUserTypeFromCookie(cookie)
	if userType != models.Admin {
		return false
	}
	return true
}
