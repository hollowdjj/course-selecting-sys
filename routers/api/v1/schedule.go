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

type CreateCourseForm struct {
	Name string `form:"name" valid:"Required;MaxSize(255)"`
	Cap  uint   `form:"cap" valid:"Required"`
}

//@Summary 创建课程
//@Produce json
//@Param name query string false "Name"
//@Param cap query uint false "Cap"
//@Success 200 {string} json "{"code":200,"data":{course_id},"msg":{"ok"}}"
//@Router /api/v1/course/create [post]
func CreateCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form CreateCourseForm
	)
	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，查看课程名称是否存在
	exist, err := models.IsCourseExistByName(form.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}
	if exist {
		appG.Response(http.StatusBadRequest, e.CourseExisted, nil)
		return
	}

	//课程不存在，那么创建课程
	id, err := models.CreateCourse(form.Name, form.Cap)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"course_id": id})
}

type GetCourseForm struct {
	CourseID uint64 `form:"course_id" valid:"Required;"`
}

//@Summary 获取课程信息
//@Produce json
//@Param course_id query uint false "CourseList"
//@Success 200 {string} json "{"code":200,"data":{course},"msg":{"ok"}}"
//@Router /api/v1/course/get [get]
func GetCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form GetCourseForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证成功后，获取课程信息
	info, err := models.GetCourseInfo(form.CourseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.CourseNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"course:": info})
}

type BindCourseForm struct {
	CourseID  uint64 `form:"course_id" valid:"Required"`
	TeacherID uint64 `form:"teacher_id" valid:"Required"`
}

//@Summary  将教师与课程绑定
//@Produce json
//@Param course_id query uint64 false "CourseList"
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/teacher/bind_course [post]
func BindCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form BindCourseForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，查询课程是否已被绑定
	bound, err := models.IsCourseBound(form.CourseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.CourseNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if bound {
		appG.Response(http.StatusBadRequest, e.CourseHasBound, nil)
		return
	}

	//若课程存在且未绑定时，绑定课程。需要注意的是，项目中的抢课部分是一个单独的算法题
	//因此，提供的teacher_id可能在数据库中根本不存在，故这里不做落库校验。
	err = models.BindCourse(form.CourseID, form.TeacherID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

//@Summary 将老师与课程解绑
//@Produce json
//@Param course_id query uint64 false "CourseList"
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/teacher/unbind_course [post]
func UnBindCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form BindCourseForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，判断课程是否被绑定过
	bound, err := models.IsCourseBound(form.CourseID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appG.Response(http.StatusBadRequest, e.CourseNotExisted, nil)
		} else {
			appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		}
		return
	}
	if !bound {
		appG.Response(http.StatusBadRequest, e.CourseNotBind, nil)
		return
	}

	//课程存在且被绑定时，解绑课程
	err = models.UnBindCourse(form.CourseID, form.TeacherID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, nil)
}

type GetTeacherCourseForm struct {
	TeacherID uint64 `form:"teacher_id" valid:"Required"`
}

//@Summary 获取一个老师的全部课程
//@Produce json
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{courses},"msg":{"ok"}}"
//@Router /api/v1/teacher/get_course [GET]
func GetTeacherCourses(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form GetTeacherCourseForm
	)

	//表单验证
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != e.OK {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//表单验证通过后，查询老师的全部课程
	courses, err := models.GetTeacherCourses(form.TeacherID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UnknownError, nil)
		return
	}

	appG.Response(http.StatusOK, e.OK, map[string]interface{}{"course_list": courses})
}

func Schedule(c *gin.Context) {
	appG := app.Gin{C: c}
	/*
		二分图匹配算法。
		例子：有六位教师：张、王、李、赵、孙、周，要安排他们去教六门课程：数学、化学、物理、语文、英语和程序设计。
		张老师会教数学、程序设计和英语；王老师会教英语和语文；李老师会教数学和物理；赵老师会教化学；孙老师会教物理
		和程序设计；周老师会教数学和物理。应该怎么样安排课程才能使每门课都有人教，每个人都只教一门课而且不至于使任
		何人去教他不懂的课？
	*/
	data := map[string][]string{
		"u1": []string{"v4", "v5"},
		"u2": []string{"v4", "v6"},
		"u3": []string{"v1", "v3"},
		"u4": []string{"v2"},
		"u5": []string{"v3", "v5"},
		"u6": []string{"v1", "v3"},
	}
	appG.Response(http.StatusOK, e.OK, utility.MaxMatch(data))
}
