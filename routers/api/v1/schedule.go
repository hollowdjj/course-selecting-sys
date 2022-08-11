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

//@Summary creat course
//@Produce json
//@Param name query string false "Name"
//@Param cap query uint false "Cap"
//@Success 200 {string} json "{"code":200,"data":{course_id},"msg":{"ok"}}"
//@Router /api/v1/course/create [post]
func CreateCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.CreateCourseForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"name": form.Name,
			"cap":  form.Cap,
			"msg":  constval.GetErrCodeMsg(errCode),
		}).Infoln("create course form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//creat course
	course := &models.Course{CourseName: form.Name, Cap: form.Cap, RemainCap: form.Cap}
	httpCode, errCode = form.CreateCourse(course)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_name": form.Name,
			"cap":         form.Cap,
			"msg":         constval.GetErrCodeMsg(errCode),
		}).Infoln("create course fail")
		appG.Response(httpCode, errCode, nil)
		return
	}

	logger.GetInstance().WithFields(logrus.Fields{
		"course_name": form.Name,
		"cap":         form.Cap,
	}).Infoln("create course succ")
	appG.Response(httpCode, errCode, map[string]interface{}{"course_id": course.CourseID})
}

//@Summary get course info
//@Produce json
//@Param course_id query uint false "CourseList"
//@Success 200 {string} json "{"code":200,"data":{course},"msg":{"ok"}}"
//@Router /api/v1/course/get [get]
func GetCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.GetCourseForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != constval.OK {
		logger.GetInstance().WithField("course_id", form.CourseID).Infoln("get course form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//get course info
	course := &models.Course{}
	httpCode, errCode = form.GetCourseInfo(course)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id": form.CourseID,
			"msg":       constval.GetErrCodeMsg(errCode),
		}).Infoln("get course info fail")
		appG.Response(httpCode, errCode, nil)
		return
	}

	logger.GetInstance().WithField("course_id", form.CourseID).Info("get course info succ")
	appG.Response(http.StatusOK, constval.OK, map[string]interface{}{"course:": *course})
}

//@Summary  bind course with teacher
//@Produce json
//@Param course_id query uint64 false "CourseList"
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/teacher/bind_course [post]
func BindCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.BindCourseForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  form.CourseID,
			"teacher_id": form.TeacherID,
		}).Infoln("bind course form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	httpCode, errCode = form.BindCourse()
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  form.CourseID,
			"teacher_id": form.TeacherID,
			"msg":        constval.GetErrCodeMsg(errCode),
		}).Infoln("bind course fail")
		appG.Response(httpCode, errCode, nil)
		return
	}

	logger.GetInstance().WithFields(logrus.Fields{
		"course_id":  form.CourseID,
		"teacher_id": form.TeacherID,
	}).Infoln("bind course succ")
	appG.Response(httpCode, errCode, nil)
}

//@Summary unbind course and teacher
//@Produce json
//@Param course_id query uint64 false "CourseList"
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{},"msg":{"ok"}}"
//@Router /api/v1/teacher/unbind_course [post]
func UnBindCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.BindCourseForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  form.CourseID,
			"teacher_id": form.TeacherID,
			"msg":        constval.GetErrCodeMsg(errCode),
		}).Infoln("unbind course form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//unbind course
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  form.CourseID,
			"teacher_id": form.TeacherID,
			"msg":        constval.GetErrCodeMsg(errCode),
		}).Infoln("unbind course fail")
		appG.Response(httpCode, errCode, nil)
		return
	}

	logger.GetInstance().WithFields(logrus.Fields{
		"course_id":  form.CourseID,
		"teacher_id": form.TeacherID,
	}).Infoln("unbind course succ")
	appG.Response(httpCode, errCode, nil)
}

//@Summary get all courses of teacher
//@Produce json
//@Param teacher_id query uint64 false "TeacherID"
//@Success 200 {string} json "{"code":200,"data":{courses},"msg":{"ok"}}"
//@Router /api/v1/teacher/get_course [GET]
func GetTeacherCourses(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.GetTeacherCourseForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, false)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"teacher_id": form.TeacherID,
			"msg":        constval.GetErrCodeMsg(errCode),
		}).Infoln("get teacher course form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//get courses
	courses := &[]models.Course{}
	httpCode, errCode = form.GetTeacherCourses(courses)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"teacher_id": form.TeacherID,
			"msg":        constval.GetErrCodeMsg(errCode),
		}).Infoln("get teacher course fail")
		appG.Response(httpCode, errCode, nil)
		return
	}

	logger.GetInstance().WithField("teacher_id", form.TeacherID).Infoln("get teacher course succ")
	appG.Response(httpCode, errCode, map[string]interface{}{"course_list": courses})
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
	appG.Response(http.StatusOK, constval.OK, utility.MaxMatch(data))
}
