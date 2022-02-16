package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"gorm.io/gorm"
	"net/http"
)

//用于验证表单
type CreateMemberForm struct {
	Username string `form:"username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `form:"password" valid:"Required;MinSize(8);MaxSize(20);PasswordCheck"`
	Nickname string `form:"nickname" valid:"Required;MinSize(4);MaxSize(20)"`
	UserType int    `form:"user_type" valid:"Required;Range(1,3)"`
}

//@Summary 创建成员。只有管理员能够访问，请求中需要带上cookie
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Param nickname query string false "Nickname"
//@Param user_type query int false "UserType"
//@Success 200 {string} json "{"code":200,"data":{"user_id"},"msg":{"ok"}}"
//@Router /api/v1/member/create [post]
func CreateMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form CreateMemberForm
	)

	//权限验证
	if !app.IsAdmin(c) {
		appG.Response(http.StatusUnauthorized, e.PermDenied, nil)
		return
	}

	//表单验证
	funcs := app.CustomFunc{
		"PasswordCheck": utility.PasswordCheck,
	}
	httpCode, errCode := app.BindAndValidCustom(c, &form, funcs)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，检查username是否已经存在
	exist, err := models.IsUserExistByName(form.Username)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	if exist {
		appG.Response(http.StatusOK, e.UserHasExisted, nil)
		return
	}

	//用户名不存在，则创建新用户
	user := models.Member{
		Username: form.Username,
		Password: form.Password,
		Nickname: form.Nickname,
		UserType: form.UserType,
	}
	if err = models.CreateMember(&user); err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"user_id:": user.UserID})
}

type UpdateMemberForm struct {
	UserID   uint64 `form:"user_id" valid:"Required"` //uint加上required，表示只接受正整数
	Nickname string `form:"nickname" valid:"Required;MinSize(4);MaxSize(20)"`
}

//@Summary 更新用户的昵称
//@Produce json
//@Param user_id query uint false "UserID"
//@Param nickname query string false "Nickname"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router	/api/v1/member/update [post]
func UpdateMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form UpdateMemberForm
	)
	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，检查username是否存在以及是否已删除
	deleted, err := models.IsUserDeleted(form.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.UserNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if deleted {
		appG.Response(http.StatusOK, e.UserHasDeleted, nil)
		return
	}

	//username存在且没有被删除，那么修改用户名
	err = models.UpdateMember(form.UserID, form.Nickname)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

type DelAndGetMemberForm struct {
	UserID uint64 `form:"user_id" valid:"Required"`
}

//@Summary  删除用户(软删除)
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/member/delete [post]
func DeleteMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DelAndGetMemberForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功，那么查询用户是否存在以及是否是已删除状态
	deleted, err := models.IsUserDeleted(form.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.UserNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if deleted {
		appG.Response(http.StatusOK, e.UserHasDeleted, nil)
		return
	}

	//用户存在且是未删除状态时，删除用户
	err = models.DeleteMember(form.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

//@Summary	获取单个成员信息
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router	/api/v1/member/ [get]
func GetMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DelAndGetMemberForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，查询用户名是否存在
	deleted, err := models.IsUserDeleted(form.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.UserNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if deleted {
		appG.Response(http.StatusBadRequest, e.UserHasDeleted, nil)
		return
	}

	//用户名存在且没有被删除，那么获取用户信息
	info, err := models.GetUserInfo(form.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"member": info})
}

//TODO offset和limit为0时显示参数错误
type GetMemberListForm struct {
	Offset int `form:"offset" valid:"Required;Min(0)"`
	Limit  int `form:"limit" valid:"Required;Min(0)"`
}

//@Summary	批量获取成员信息
//@Produce json
//@Param offset query uint false "OffSet"
//@Param limit query uint false "Limit"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router /api/v1/member/list [get]
func GetMembers(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form GetMemberListForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//批量获取信息
	members, err := models.GetUserInfoList(form.Offset, form.Limit)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"member_list": members})
}
