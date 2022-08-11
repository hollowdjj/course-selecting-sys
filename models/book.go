package models

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	//group name: student_course. There are two types of key-val pair:
	//1. key is format as userid_courseid and val is "1", which is used to judge whether user has this course
	//2. key is userid and val is a string which is format as courseid1_courseid2_....., which is used to get student course
	maxStudentCourseBytes int64 = 128 * 1024 * 1024
	//group name: course_remain_cap		 key is courseid and val is remain cap
	maxCourseRemainCapBytes int64 = 10 * 1024 * 1024
)

//used for booking course
type BookCourseForm struct {
	UserID   string `valid:"Required;Numeric"`
	CourseID string `valid:"Required;Numeric"`
}

func (b BookCourseForm) BookCourse() (int, constval.ErrNo) {
	//check whether student has this course
	studentCourseCache := cache.GetGroupCache("student_course")
	if studentCourseCache == nil {
		studentCourseCache = cache.NewGroupCache("student_course", maxStudentCourseBytes, cache.GetterFunc(StudentCourseGetter))
	}
	key := b.UserID + "_" + b.CourseID
	val, err := studentCourseCache.Get(key, cache.DefaultOption)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"form": b,
			"err":  err,
		}).Errorln("get student course info err")
	}
	if val.Len() > 0 {
		return http.StatusOK, constval.StudentHasCourse
	}

	//get course remain cap and judge
	courseRemainCapCache := cache.GetGroupCache("course_remain_cap")
	if courseRemainCapCache == nil {
		courseRemainCapCache = cache.NewGroupCache("course_remain_cap", maxCourseRemainCapBytes, cache.GetterFunc(CourseRemainCapGetter))
	}
	val, err = courseRemainCapCache.Get(b.CourseID, cache.DefaultOption)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"userid":   b.UserID,
			"courseid": b.CourseID,
			"err":      err,
		}).Errorln("get course remain cap err")
		return http.StatusInternalServerError, constval.UnknownError
	}
	remainCap, err := strconv.Atoi(string(val.ByteSlice()))
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"data": val.ByteSlice(),
			"err":  err,
		}).Errorln("convert course remain_cap to int error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if remainCap <= 0 {
		return http.StatusOK, constval.CourseNotAvailable
	}

	//book course and update cache
	remainCap--
	courseRemainCapCache.Add(key, []byte(strconv.Itoa(remainCap)), 60)
	result := Db.Model(&Course{}).Where("course_id = ? AND remain_cap > 0", b.CourseID).
		Update("remain_cap", gorm.Expr("remain_cap - ?", 1))
	if result.Error != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id": b.CourseID,
			"err":       err,
		}).Errorln("update course remain_cap error")
		courseRemainCapCache.Add(key, []byte(strconv.Itoa(remainCap+1)), 60)
		return http.StatusInternalServerError, constval.UnknownError
	}
	httpCode, errCode := http.StatusOK, constval.OK
	if result.RowsAffected == 0 {
		//suggest that course has no cap and cache is not up to date
		remainCap = 0
		httpCode, errCode = http.StatusOK, constval.CourseNotAvailable
		courseRemainCapCache.Add(key, []byte("0"), 60)
	}

	return httpCode, errCode
}

//used for quering student course
type GetStudentCourseForm struct {
	UserID string `valid:"Required;Numeric"`
}

func (g GetStudentCourseForm) GetStudentCourse(courses *[]Course) (int, constval.ErrNo) {
	studentCourseCache := cache.GetGroupCache("student_course")
	if studentCourseCache == nil {
		studentCourseCache = cache.NewGroupCache("student_course", maxStudentCourseBytes, cache.GetterFunc(StudentCourseGetter))
	}
	val, err := studentCourseCache.Get(g.UserID, cache.DefaultOption)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"user_id": g.UserID,
			"err":     err,
		}).Errorln("get student courses error")
		return http.StatusInternalServerError, constval.UnknownError
	}

	courseInfoCache := cache.GetGroupCache("course_info")
	if courseInfoCache == nil {
		courseInfoCache = cache.NewGroupCache("course_info", maxCourseInfoCacheBytes, cache.GetterFunc(CourseInfoGetter))
	}
	coursesID := strings.Split(val.String(), "_")
	for i := range coursesID {
		val, err := courseInfoCache.Get(coursesID[i], cache.DefaultOption)
		if err != nil {
			logger.GetInstance().WithFields(logrus.Fields{
				"course_id": coursesID,
				"err":       err,
			}).Errorln("get course info error")
			return http.StatusInternalServerError, constval.UnknownError
		}
		//unmarshal
		course := Course{}
		err = json.Unmarshal(val.ByteSlice(), &course)
		if err != nil {
			logger.GetInstance().WithField("err", err).Errorln("json unmarshal course info error")
			continue
		}
		*courses = append(*courses, course)
	}
	return http.StatusOK, constval.OK
}
