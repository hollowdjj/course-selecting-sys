package models

import (
	"encoding/json"
	"net/http"

	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"github.com/sirupsen/logrus"
)

const (
	//the length of a token is 128 bits(md5). In group cache token, there are
	//key-value pairs: token:userinfo, 128 * 1024 * 1024 is already big enough.
	loginCacheMaxBytes = 128 * 1024 * 1024
	loginCacheTTL      = 3 * 3600
)

type LoginForm struct {
	Username string `form:"username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `form:"password" valid:"Required;MinSize(8);MaxSize(20)"`
}

//query db to check if username and password match
func (l LoginForm) Login(tok *string) (int, constval.ErrNo) {
	groupCacheLogin := cache.GetGroupCache("login")
	if groupCacheLogin == nil {
		groupCacheLogin = cache.NewGroupCache("login", loginCacheMaxBytes, cache.GetterFunc(UserInfoGetter))
	}
	_, err := groupCacheLogin.Get(l.Username, cache.Option{
		FromLocal:  true,
		FromPeer:   false,
		FromGetter: true,
		TTL:        loginCacheTTL,
	})
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"username": l.Username,
			"err":      err,
		}).Errorln("user login error")
	}

	//login succ, generate token and add cache
	token := utility.GenerateToken()
	groupCacheLogin.Add(token, []byte(l.Username), loginCacheTTL)
	tok = &token
	return http.StatusOK, constval.OK
}

//judge is user is admin according to token
func IsAdmin(token string) (int, constval.ErrNo) {
	groupCacheToken := cache.GetGroupCache("login")
	if groupCacheToken == nil {
		return http.StatusBadRequest, constval.LoginRequired
	}
	val, _ := groupCacheToken.Get(token, cache.Option{
		FromLocal:  true,
		FromPeer:   false,
		FromGetter: false,
	})
	if val.Len() == 0 {
		return http.StatusUnauthorized, constval.LoginRequired
	}
	userinfo := &UserInfo{}
	err := json.Unmarshal(val.ByteSlice(), userinfo)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("json unmarshal fail")
		return http.StatusInternalServerError, constval.UnknownError
	}
	return http.StatusOK, constval.OK
}
