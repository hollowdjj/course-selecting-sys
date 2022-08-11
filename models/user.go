package models

import (
	"encoding/json"
	"net/http"

	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	maxUserInfoCacheBytes int64 = 64 * 1024 * 1024
)

//used for creating user
type CreateUserForm struct {
	Username string `json:"username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `json:"password" valid:"Required;MinSize(8);MaxSize(20);PasswordCheck"`
	Nickname string `json:"nickname" valid:"Required;MinSize(4);MaxSize(20)"`
	UserType int    `json:"user_type" valid:"Required;Range(1,3)"`
}

func (c *CreateUserForm) CreateUser(user *User) (int, constval.ErrNo) {
	result := Db.FirstOrCreate(user, *user)
	if err := result.Error; err != nil {
		logger.GetInstance().WithField("err", err).Errorln("create user error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if result.RowsAffected == 0 {
		logger.GetInstance().WithField("username", c.Username).Infoln("username already exist")
		return http.StatusBadRequest, constval.UserExisted
	}

	return http.StatusOK, constval.OK
}

//used for update user's nickname
type UpdateMemberForm struct {
	UserID   uint64 `json:"user_id" valid:"Required"` //uint加上required，表示只接受正整数
	Nickname string `json:"nickname" valid:"Required;MinSize(4);MaxSize(20)"`
}

func (u *UpdateMemberForm) UpdateUser() (int, constval.ErrNo) {
	result := Db.Model(&User{}).Where("user_id = ?", u.UserID).Update("nickname", u.Nickname)
	if result.Error != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"user_id":  u.UserID,
			"nickname": u.Nickname,
		}).Errorln("update user nickname error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	return http.StatusOK, constval.OK
}

//used for del or get user
type DelOrGetUserForm struct {
	UserID string `form:"user_id" json:"user_id" valid:"Required"`
}

func (d *DelOrGetUserForm) DeleteUser() (int, constval.ErrNo) {
	result := Db.Model(&User{}).Where("user_id = ?", d.UserID).Update("is_active", 0)
	if result.Error != nil {
		logger.GetInstance().WithField("user_id", d.UserID).Errorln("del user error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if result.RowsAffected == 0 {
		logger.GetInstance().WithField("user_id", d.UserID).Infoln("user deleted or not exist")
		return http.StatusBadRequest, constval.UserDeletedOrNotExist
	}
	return http.StatusOK, constval.OK
}

func (d *DelOrGetUserForm) GetUserInfo(userInfo *UserInfo) (int, constval.ErrNo) {
	groupCacheUser := cache.GetGroupCache("user")
	if groupCacheUser == nil {
		groupCacheUser = cache.NewGroupCache("user", maxUserInfoCacheBytes, cache.GetterFunc(UserInfoGetterByUserID))
	}
	val, err := groupCacheUser.Get(d.UserID, cache.DefaultOption)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"group":  "user",
			"userid": d.UserID,
		}).Errorln("get user info error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	//unmarshal user info
	err = json.Unmarshal(val.ByteSlice(), userInfo)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("json unmarshal error")
		return http.StatusInternalServerError, constval.UnknownError
	}

	return http.StatusOK, constval.OK
}

//used from get user list
type GetUserListForm struct {
	Offset int `form:"offset" valid:"Min(0)"`
	Limit  int `form:"limit" valid:"Min(-1)"`
}

func (g *GetUserListForm) GetUserList(userList *[]UserInfo) (int, constval.ErrNo) {
	err := Db.Model(&User{}).Offset(g.Offset).Limit(g.Limit).Find(userList).Error
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"offset": g.Offset,
			"limit":  g.Limit,
			"err":    err,
		}).Errorln("get user list error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	return http.StatusOK, constval.OK
}
