package auth_service

import (
	"github.com/go-redis/redis"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/gredis"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"time"
)

type Auth struct {
	Username string `form:"username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `form:"password" valid:"Required;MinSize(8);MaxSize(20)"`
}

//在redis中保存username:token token:userinfo
func (a *Auth) Login() (bool, string, error) {
	//查数据库验证用户名和密码是否正确
	pass, member, err := models.CheckAuth(a.Username, a.Password)
	if !pass || err != nil {
		logging.Error(err)
		return false, "", err
	}

	//删除旧的token如果存在的话
	oldToken, err := gredis.Rdb.Get(a.Username).Result()
	if err != redis.Nil {
		gredis.Delete(oldToken)
	}

	//生成新的token
	token := utility.GenerateToken()
	err = gredis.Set(a.Username, token, 10*time.Second)
	if err != nil {
		logging.Error(err)
		return false, "", err
	}
	err = gredis.Set(token, member, 10*time.Second)
	if err != nil {
		logging.Error(err)
		return false, "", err
	}

	return true, token, nil
}
