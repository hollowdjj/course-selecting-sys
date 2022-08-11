package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/models"
	"github.com/hollowdjj/course-selecting-sys/pkg/app"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
)

func BookCourse(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.BookCourseForm
	)

	//form validation
	httpCode, errCode := app.BindAndValid(c, &form, true)
	if errCode != constval.OK {
		logger.GetInstance().WithFields(logrus.Fields{
			"userid":   form.UserID,
			"courseid": form.CourseID,
			"msg":      constval.GetErrCodeMsg(errCode),
		}).Infoln("book course form invalid")
		appG.Response(httpCode, errCode, nil)
		return
	}

	//book course
	httpCode, errCode = form.BookCourse()
}

func GetStudentCourse(c *gin.Context) {

}
