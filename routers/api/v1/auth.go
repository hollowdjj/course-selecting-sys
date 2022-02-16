package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/setting"
	"github.com/hollowdjj/course-selecting-sys/pkg/utility"
	"net/http"
)

//用于登录时的表单验证。form标签用于让gin验证URL中是否带有username以及password参数
//valid标签用于表单参数是否合法的验证。
type AuthForm struct {
	Username string `form:"username" valid:"Required;MinSize(8);MaxSize(20)"`
	Password string `form:"password" valid:"Required;MinSize(8);MaxSize(20)"`
}

//@Summary 用户登录
//@Produce json
//@Param username query string false "UserName"
//@Param password query string false "Password"
//@Success 200 {string} json "{"code":200,"data":{UserID},"msg":{"ok"}}"
//@Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AuthForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过，检查用户名和密码是否正确
	pass, err, value := models.CheckAuth(form.Username, form.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	if !pass {
		appG.Response(http.StatusUnauthorized, e.WrongPassword, nil)
		return
	}

	//登录成功，那么设置cookie
	//设置cookie。如果请求中没有cookie，err为ErrNoCookie。cookie的值设置为用户ID+用户名+用户类型
	_, err = appG.C.Cookie("camp-seesion")
	if err != nil {
		//请求中没有cookie，那么设置cookie
		cookieValue := fmt.Sprintf("%d+%s+%v", value.UserID, form.Username, value.UserType)
		appG.C.SetCookie("camp-seesion", cookieValue,
			1800, "/", setting.AppSetting.Host, false, true)
	}

	//登录成功时需要返回登录用户的ID
	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"user_id": value.UserID})
}

//Logout 登出
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}
	//清除cookie。这里设置0或-1都不行，不知道原因是什么
	cookie, err := appG.C.Cookie("camp-seesion")
	if err == http.ErrNoCookie {
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		return
	}
	appG.C.SetCookie("camp-seesion", cookie, 1,
		"/", setting.AppSetting.Host, false, true)
	appG.Response(http.StatusOK, e.OK, nil)
}

//WhoAmI 获取个人信息
func WhoAmI(c *gin.Context) {
	appG := app.Gin{C: c}

	//从cookie中得到cookie值。cookie值即为用户ID+用户名+type
	cookie, err := c.Cookie("camp-seesion")
	if err == http.ErrNoCookie {
		appG.Response(http.StatusUnauthorized, e.LoginRequired, nil)
		return
	}

	//查询用户信息并响应
	id := utility.GetUserIDFromCookie(cookie)
	userInfo, err := models.GetUserInfo(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	appG.Response(http.StatusOK, e.OK, userInfo)
}
