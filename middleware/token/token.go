package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
)

/*
Gin中间件的作用有两个：
1. Web请求到达我们定义的HTTP请求处理方法之前，拦截请求并进行相应处理(比如：权限验证，数据过滤等)
2. 在我们处理完请求并响应客户端时，拦截响应并进行相应的处理(比如：添加统一响应头或数据格式等)
在Gin框架中，中间件就是一个函数，其函数类型为type HandlerFunc func(*gin.Context)，就是一个参数类型
为*gin.Context且没有返回值的函数
*/

//a token middleware
func Token(c *gin.Context) {
	appG := app.Gin{C: c}

	//get token from request header
	token := c.GetHeader("Authorization")
	if token == "" {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		c.Abort()
		return
	}

	//find token in local cache
	groupCacheLogin := cache.GetGroupCache("login")
	if groupCacheLogin == nil {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		c.Abort()
		return
	}
	val, err := groupCacheLogin.Get(token, cache.Option{
		FromLocal:  true,
		FromPeer:   false,
		FromGetter: false,
	})
	if err != nil {
		appG.Response(http.StatusInternalServerError, constval.UnknownError, nil)
		c.Abort()
		return
	}
	if val.Len() == 0 {
		appG.Response(http.StatusUnauthorized, constval.LoginRequired, nil)
		c.Abort()
		return
	}

	//go next if find token in local cache
	c.Next()
}
