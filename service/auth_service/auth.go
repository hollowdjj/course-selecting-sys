package auth_service

import (
	"github.com/go-redis/redis"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/gredis"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"github.com/hollowdjj/course-selecting-sys/service/memeber_service"
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
	err = gredis.Rdb.Set(a.Username, token, 24*time.Hour).Err()
	if err != nil {
		logging.Error(err)
		return false, "", err
	}
	err = gredis.Set(token, member, 24*time.Hour)
	if err != nil {
		logging.Error(err)
		return false, "", err
	}

	return true, token, nil
}

//根据token判断用户是否是管理员
func IsAdmin(token string) (int, e.ErrNo) {
	var member models.MemberInfo
	httpCode, errCode := memeber_service.GetMemberInfoFromRedis(token, &member)
	if member.UserType != 1 {
		errCode = e.PermDenied
	}

	return httpCode, errCode
}
