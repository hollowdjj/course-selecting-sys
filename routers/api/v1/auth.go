package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/gredis"
	"github.com/hollowdjj/course-selecting-sys/service/auth_service"
	"net/http"
)

//@Summary 用户登录
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Success 200 {string} json "{"code":200,"data":{UserID},"msg":{"ok"}}"
//@Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		auth auth_service.Auth
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &auth, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过，那么登录
	pass, token, err := auth.Login()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	if !pass {
		appG.Response(http.StatusUnauthorized, e.WrongPassword, nil)
		return
	}

	//登录成功返回登录用户的ID以及服务端生成的token
	c.Header("Authorization", token)
	appG.Response(http.StatusOK, e.OK, nil)
}

//Logout 登出
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}

	//获取用户的token，并根据token在redis中获取用户信息
	token := c.GetHeader("Authorization")
	mem, err := gredis.Get(token)
	if err == redis.Nil {
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		return
	}

	var memberAuth models.MemberInfo
	err = json.Unmarshal(mem, &memberAuth)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	//删除token
	_, err = gredis.Delete(token)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	_, err = gredis.Delete(memberAuth.Username)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

//WhoAmI 获取个人信息
func WhoAmI(c *gin.Context) {
	appG := app.Gin{C: c}

	//获取token
	token := c.GetHeader("Authorization")

	//查询redis得到用户信息
	mem, err := gredis.Get(token)
	if err == redis.Nil {
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		return
	}
	var member models.MemberInfo
	err = json.Unmarshal(mem, &member)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"member_info": member})
}
