package memeber_service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/gredis"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"gorm.io/gorm"
	"net/http"
	"time"
)

//创建成员
type CreateMemberForm struct {
	Username string `json:"username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `json:"password" valid:"Required;MinSize(8);MaxSize(20);PasswordCheck"`
	Nickname string `json:"nickname" valid:"Required;MinSize(4);MaxSize(20)"`
	UserType int    `json:"user_type" valid:"Required;Range(1,3)"`
}

func (c *CreateMemberForm) CreateMember(member *models.Member) (int, e.ErrNo) {
	//先查redis判断该用户在短时间之前是否已经创建过了
	exist, _ := gredis.Exist("create:" + c.Username)
	if exist {
		return http.StatusBadRequest, e.UserHasExisted
	}

	//开启一个事务创建用户
	err := models.Db.Transaction(func(tx *gorm.DB) error {
		//首先判断用户名是否存在
		exist, err := models.IsUserExistByName(c.Username)
		if err != nil {
			//回滚
			logging.Error(err)
			return err
		}

		//用户不存在那么创建用户
		if !exist {
			member.Username = c.Username
			member.Password = c.Password
			member.Nickname = c.Nickname
			member.UserType = c.UserType
			err = models.CreateMember(member)
			if err != nil {
				//回滚
				logging.Error(err)
				return err
			}
		}

		//返回nil提交事务
		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, e.UnknownError
	}

	//添加一个10秒的redis缓存，避免短时间内多个相同的username
	gredis.Rdb.Set("create:"+c.Username, 1, 10*time.Second)

	if member.UserID > 0 {
		return http.StatusOK, e.OK
	}

	return http.StatusBadRequest, e.UserHasExisted
}

//更新成员
type UpdateMemberForm struct {
	UserID   uint64 `json:"user_id" valid:"Required"` //uint加上required，表示只接受正整数
	Nickname string `json:"nickname" valid:"Required;MinSize(4);MaxSize(20)"`
}

func (u *UpdateMemberForm) UpdateMember() (int, e.ErrNo) {
	httpCode, errCode := http.StatusOK, e.OK
	_ = models.Db.Transaction(func(tx *gorm.DB) error {
		//使用当前读判断用户是否存在或者是否已删除。
		deleted, err := models.IsUserDeleted(u.UserID, true)
		if err == gorm.ErrRecordNotFound {
			httpCode, errCode = http.StatusBadRequest, e.UserNotExisted
			return nil
		}
		if deleted {
			httpCode, errCode = http.StatusBadRequest, e.UserHasDeleted
			return nil
		}

		//更新用户
		err = models.UpdateMember(u.UserID, u.Nickname)
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.UnknownError
			//回滚
			return err
		}

		//提交事务
		return nil
	})

	return httpCode, errCode
}

//删除或获取成员
type DelAndGetMemberForm struct {
	UserID uint64 `form:"user_id" json:"user_id" valid:"Required"`
}

type MemberInfoCache struct {
	Info    models.MemberInfo
	Deleted bool
	Existed bool
}

func (d *DelAndGetMemberForm) DeleteMember() (int, e.ErrNo) {
	httpCode, errCode := http.StatusOK, e.OK
	_ = models.Db.Transaction(func(tx *gorm.DB) error {
		//使用当前读判断用户是否存在或者是否已删除。
		deleted, err := models.IsUserDeleted(d.UserID, true)
		if err == gorm.ErrRecordNotFound {
			httpCode, errCode = http.StatusBadRequest, e.UserNotExisted
			return nil
		}
		if deleted {
			httpCode, errCode = http.StatusBadRequest, e.UserHasDeleted
			return nil
		}

		//删除用户
		err = models.DeleteMember(d.UserID)
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.UnknownError
			//回滚
			return err
		}

		return nil
	})

	return httpCode, errCode
}

func (d *DelAndGetMemberForm) GetMemberInfo(memberInfo *models.MemberInfo) (int, e.ErrNo) {
	//首先查redis缓存
	memberCache := MemberInfoCache{Info: *memberInfo}
	mem, err := gredis.Get(fmt.Sprintf("%d", d.UserID))
	if err != nil {
		if err != redis.Nil {
			logging.Error(err)
			return http.StatusInternalServerError, e.UnknownError
		}
	} else {
		//缓存中有就直接返回
		err = json.Unmarshal(mem, &memberCache)
		if err != nil {
			logging.Error(err)
			return http.StatusInternalServerError, e.UnknownError
		}

		if !memberCache.Existed {
			return http.StatusBadRequest, e.UserNotExisted
		}
		if memberCache.Deleted {
			return http.StatusBadRequest, e.UserHasDeleted
		}

		*memberInfo = memberCache.Info
		return http.StatusOK, e.OK
	}

	//redis中没有就再查数据库
	httpCode, errCode := http.StatusOK, e.OK

	err = models.Db.Transaction(func(tx *gorm.DB) error {
		//使用当前读判断用户是否存在或者是否已删除。
		deleted, err := models.IsUserDeleted(d.UserID, true)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				httpCode, errCode = http.StatusBadRequest, e.UserNotExisted
				memberCache.Existed = false
				return nil
			}
			return err
		}

		if deleted {
			httpCode, errCode = http.StatusBadRequest, e.UserHasDeleted
			memberCache.Deleted = true
			return nil
		}

		//查询用户信息
		info, err := models.GetUserInfo(d.UserID)
		if err != nil {
			return err
		}

		*memberInfo = info
		memberCache.Info = *memberInfo
		memberCache.Existed = true
		return nil
	})
	if err != nil {
		logging.Error(err)
		httpCode, errCode = http.StatusInternalServerError, e.UnknownError
	}

	//添加一个redis缓存
	_ = gredis.Set(fmt.Sprintf("%d", d.UserID), memberCache, 10*time.Second)

	return httpCode, errCode
}

//批量获取成员信息
type GetMemberListForm struct {
	Offset int `form:"offset" valid:"Min(0)"`
	Limit  int `form:"limit" valid:"Min(-1)"`
}

func (g *GetMemberListForm) GetMemberList() ([]models.MemberInfo, int, e.ErrNo) {
	members, err := models.GetUserInfoList(g.Offset, g.Limit)
	if err != nil {
		logging.Error(err)
		return nil, http.StatusInternalServerError, e.UnknownError
	}
	return members, http.StatusOK, e.OK
}

//根据token在redis中查找用户信息。若token不存在，返回e.LoginRequired。
func GetMemberInfoFromRedis(token string, memberInfo *models.MemberInfo) (int, e.ErrNo) {
	mem, err := gredis.Get(token)
	if err == redis.Nil {
		//token不存在，未登录。
		return http.StatusUnauthorized, e.LoginRequired
	}

	//解码用户信息
	err = json.Unmarshal(mem, memberInfo)
	if err != nil {
		return http.StatusInternalServerError, e.UnknownError
	}

	return http.StatusOK, e.OK
}
