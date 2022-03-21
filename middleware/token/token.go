package token

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/gredis"
	"github.com/hollowdjj/course-selecting-sys/service/memeber_service"
	"net/http"
)

/*
Gin中间件的作用有两个：
1. Web请求到达我们定义的HTTP请求处理方法之前，拦截请求并进行相应处理(比如：权限验证，数据过滤等)
2. 在我们处理完请求并响应客户端时，拦截响应并进行相应的处理(比如：添加统一响应头或数据格式等)
在Gin框架中，中间件就是一个函数，其函数类型为type HandlerFunc func(*gin.Context)，就是一个参数类型
为*gin.Context且没有返回值的函数
*/

//创建一个gin中间件
func Token(c *gin.Context) {
	appG := app.Gin{C: c}

	//获取header中的token
	token := c.GetHeader("Authorization")
	if token == "" {
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		c.Abort()
		return
	}

	//在redis中查找token是否存在
	exist, err := gredis.Exist(token)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		c.Abort()
		return
	}
	if !exist {
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		c.Abort()
		return
	}

	//token存在，那么在redis中由token找到username，然后再对比username对应的token
	var member models.MemberInfo
	httpCode, errCode := memeber_service.GetMemberInfoFromRedis(token, &member)
	if errCode != e.OK {
		gredis.Delete(token)
		appG.Response(httpCode, errCode, nil)
		c.Abort()
		return
	}

	storedToken, err := gredis.Rdb.Get(member.Username).Result()
	if err == redis.Nil || storedToken != token {
		gredis.Delete(token)
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		c.Abort()
		return
	}

	c.Next()
}
