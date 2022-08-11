package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/hollowdjj/course-selecting-sys/routers/api/v1"
)

func RegisterRouter() *gin.Engine {
	//新建一个gin路由并绑定中间件
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery())

	//设置路由
	apiv1 := g.Group("/api/v1")
	{
		//switch work mode
		apiv1.POST("/proxy/switch")

		apiv1.POST("/auth/login", v1.Login)   //登录
		apiv1.POST("/auth/logout", v1.Logout) //登出
		apiv1.GET("/auth/whoami", v1.WhoAmI)  //获取个人信息

		//成员
		apiv1.POST("/member/create", v1.CreateUser)
		apiv1.GET("/member/", v1.GetUser)
		apiv1.GET("/member/list", v1.GetUsers)
		apiv1.POST("/member/update", v1.UpdateUser)
		apiv1.POST("/member/delete", v1.DeleteUser)

		//排课
		apiv1.POST("/course/create", v1.CreateCourse)
		apiv1.GET("/course/get", v1.GetCourse)
		apiv1.POST("/teacher/bind_course", v1.BindCourse)
		apiv1.POST("/teacher/unbind_course", v1.UnBindCourse)
		apiv1.GET("/teacher/get_course", v1.GetTeacherCourses)
		apiv1.POST("/course/schedule", v1.Schedule)

		//抢课
		apiv1.POST("/student/book_course", v1.BookCourse)
		apiv1.GET("/student/course", v1.GetStudentCourse)
	}

	return g
}
