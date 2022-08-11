package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
)

//@Summary user login
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Success 200 {string} json "{"code":200,"data":{UserID},"msg":{"ok"}}"
//@Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.LoginForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"username": form.Username,
			"msg":      constval.GetErrCodeMsg(errCode),
		}).Infoln("login form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//try login
	var token *string
	httpCode, errCode = form.Login(token)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"username": form.Username,
			"msg":      constval.GetErrCodeMsg(errCode),
		}).Infoln("login fail")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//response token to client
	logger.GetInstance().WithField("user", form.Username).Infoln("user login succ")
	c.Header("Authorization", *token)
	appG.Response(httpCode, errCode, nil)
}

//Logout
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}

	//get token
	token := c.GetHeader("Authorization")

	//del cache
	groupCacheToken := cache.GetGroupCache("login")
	if groupCacheToken == nil {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		return
	}
	val, _ := groupCacheToken.Get(token, cache.Option{
		FromLocal:  true,
		FromPeer:   false,
		FromGetter: false,
	})
	if val.Len() == 0 {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		return
	}

	//unmarshal user info to get username
	userInfo := &models.UserInfo{}
	err := json.Unmarshal(val.ByteSlice(), userInfo)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("json unmarshal fail")
	}
	//del cache
	groupCacheToken.Del(token)
	groupCacheToken.Del(userInfo.Username)
	appG.Response(http.StatusOK, constval.OK, nil)
	logger.GetInstance().WithField("user", userInfo.Username).Infoln("user logout succ")
}

//WhoAmI get user infomation
func WhoAmI(c *gin.Context) {
	appG := app.Gin{C: c}

	//get token
	token := c.GetHeader("Authorization")

	//look up cache
	groupCacheToken := cache.GetGroupCache("login")
	if groupCacheToken == nil {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		return
	}
	val, _ := groupCacheToken.Get(token, cache.Option{
		FromLocal:  true,
		FromPeer:   false,
		FromGetter: false,
	})
	if val.Len() == 0 {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		return
	}

	//unmarshal user info
	userInfo := &models.UserInfo{}
	err := json.Unmarshal(val.ByteSlice(), userInfo)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("json unmarshal fail")
		appG.Response(http.StatusInternalServerError, constval.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, constval.OK, map[string]interface{}{"user_info": userInfo})
}
