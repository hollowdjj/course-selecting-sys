package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"github.com/hollowdjj/course-selecting-sys/service/auth_service"
	"github.com/hollowdjj/course-selecting-sys/service/memeber_service"
	"net/http"
)

//@Summary 创建成员，只有管理员能够访问。
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
		form memeber_service.CreateMemberForm
	)

	//获取token并权限验证
	token := c.GetHeader("Authorization")
	httpCode, errCode := auth_service.IsAdmin(token)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证
	funcs := app.CustomFunc{
		"PasswordCheck": utility.PasswordCheck,
	}
	httpCode, errCode = app.BindAndValidCustom(c, &form, funcs, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，创建成员
	var member models.Member
	httpCode, errCode = form.CreateMember(&member)
	appG.Response(httpCode, errCode, map[string]interface{}{"user_id:": member.UserID})
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
		form memeber_service.UpdateMemberForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，更新用户
	httpCode, errCode = form.UpdateMember()
	appG.Response(httpCode, errCode, nil)
}

//@Summary  删除用户(软删除)
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/member/delete [post]
func DeleteMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form memeber_service.DelAndGetMemberForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功，那么查询用户信息
	httpCode, errCode = form.DeleteMember()

	appG.Response(httpCode, errCode, nil)
}

//@Summary	获取单个成员信息
//@Produce json
//@Param user_id query uint false "UserID"
//@Success 200 {string} json "{"code":200,"data":{userinfo},"msg":{"ok"}}"
//@Router	/api/v1/member/ [get]
func GetMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form memeber_service.DelAndGetMemberForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，查询用户信息
	var memberInfo models.MemberInfo
	httpCode, errCode = form.GetMemberInfo(&memberInfo)

	appG.Response(httpCode, errCode, map[string]interface{}{"member": memberInfo})
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
		form memeber_service.GetMemberListForm
	)

	//表单验证
	_, b1 := c.GetQuery("offset")
	_, b2 := c.GetQuery("limit")
	if !b1 || !b2 {
		appG.Response(http.StatusOK, e.ParamInvalid, nil)
		return
	}
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//批量获取信息
	members, httpCode, errCode := form.GetMemberList()

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"member_list": members})
}
