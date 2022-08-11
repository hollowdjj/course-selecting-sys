package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"github.com/sirupsen/logrus"
)

//@Summary create member
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Param nickname query string false "Nickname"
//@Param user_type query int false "UserType"
//@Success 200 {string} json "{"code":200,"data":{"user_id"},"msg":{"ok"}}"
//@Router /api/v1/member/create [post]
func CreateUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.CreateUserForm
	)
	c.
	//check if user is admin
	token := c.GetHeader("Authorization")
	httpCode, errCode := models.IsAdmin(token)
	if errCode != constval.OK {
		logger.GetInstance().Errorln("user is not admin")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//form validation
	funcs := app.CustomFunc{
		"PasswordCheck": utility.PasswordCheck,
	}
	httpCode, errCode = app.BindAndValidCustom(c, &form, funcs, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"username": form.Username,
			"password": form.Password,
			"usertype": form.UserType,
			"nickname": form.Nickname,
		}).Infoln("create user request form incorrect")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//creat member
	user := &models.User{
		Username: form.Username,
		Password: form.Password,
		UserType: form.UserType,
		Nickname: form.Nickname,
	}
	httpCode, errCode = form.CreateUser(user)
	appG.Response(httpCode, errCode, map[string]interface{}{"user_id:": user.UserID})
	logger.GetInstance().WithFields(logrus.Fields{
		"username": form.Username,
		"usertype": form.UserType,
		"nickname": form.Nickname,
	}).Infoln("create user succ")
}

//@Summary update user's nickname
//@Produce json
//@Param user_id query uint false "UserID"
//@Param nickname query string false "Nickname"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router	/api/v1/member/update [post]
func UpdateUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.UpdateMemberForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"user_id":   form.UserID,
			"nick_name": form.Nickname,
		}).Infoln("update user request form incorrect")
		appG.Response(httpCode, errCode, nil)
		return
	}
	
	//update user
	httpCode, errCode = form.UpdateUser()
	msg := ""
	if errCode == constval.OK {
		msg = "update user nickname succ"
	} else {
		msg = "update user nickname fail"
	}
	logger.GetInstance().WithFields(logrus.Fields{
		"user_id":   form.UserID,
		"nick_name": form.Nickname,
	}).Errorln(msg)

	appG.Response(httpCode, errCode, nil)
}

//@Summary  删除用户(软删除)
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/member/delete [post]
func DeleteUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.DelOrGetUserForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithField("user_id", form.UserID).Infoln("del user form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//delete user
	httpCode, errCode = form.DeleteUser()
	appG.Response(httpCode, errCode, nil)
	msg := ""
	if errCode == constval.OK {
		msg = "del user succ"
	} else {
		msg = "del user fail"
	}
	logger.GetInstance().WithField("user_id", form.UserID).Infoln(msg)
}

//@Summary	获取单个成员信息
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router	/api/v1/member/ [get]
func GetUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.DelOrGetUserForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != constval.OK {
		logger.GetInstance().WithField("user_id", form.UserID).Infoln("get user form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//get user info
	var userInfo models.UserInfo
	httpCode, errCode = form.GetUserInfo(&userInfo)
	msg := ""
	if errCode == constval.OK {
		msg = "get user info succ"
	} else {
		msg = "get user info fail"
	}
	logger.GetInstance().WithField("user_id", form.UserID).Info(msg)
	appG.Response(httpCode, errCode, map[string]interface{}{"user": userInfo})
}

//@Summary	批量获取成员信息
//@Produce json
//@Param offset query uint false "OffSet"
//@Param limit query uint false "Limit"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router /api/v1/member/list [get]
func GetUsers(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.GetUserListForm
	)

	//form validation
	offset, b1 := c.GetQuery("offset")
	limit, b2 := c.GetQuery("limit")
	if !b1 || !b2 {
		logger.GetInstance().WithFields(logrus.Fields{
			"offset": offset,
			"limit":  limit,
		}).Infoln("get user list form invalid")
		appG.Response(http.StatusBadGateway, constval.ParamInvalid, nil)
		return
	}
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"offset": offset,
			"limit":  limit,
		}).Infoln("get users form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//get user list
	userList := make([]models.UserInfo, form.Limit)
	httpCode, errCode = form.GetUserList(&userList)
	if errCode != constval.OK {
		logger.GetInstance().Errorln("get user list fail")
	}

	logger.GetInstance().Errorln("get user list succ")
	appG.Response(http.StatusOK, constval.OK, map[string]interface{}{"user_list": userList})
}
