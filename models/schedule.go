package models

import (
	"encoding/json"
	"net/http"

	"github.com/hollowdjj/course-selecting-sys/cache"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	maxCourseInfoCacheBytes = 128 * 1024 * 1024 //128MB
)

//used for creating course
type CreateCourseForm struct {
	Name string `form:"name" valid:"Required;MaxSize(255)"`
	Cap  uint   `form:"cap" valid:"Required"`
}

//create course if not exist
func (c CreateCourseForm) CreateCourse(course *Course) (int, constval.ErrNo) {
	result := Db.FirstOrCreate(course)
	if err := result.Error; err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"name": c.Name,
			"cap":  c.Cap,
			"err":  err,
		}).Errorln("create course error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if result.RowsAffected == 0 {
		logger.GetInstance().WithFields(logrus.Fields{
			"name": c.Name,
			"cap":  c.Cap,
		}).Errorln("course already exist")
		return http.StatusBadRequest, constval.CourseExisted
	}
	return http.StatusOK, constval.OK
}

//used for getting course info
type GetCourseForm struct {
	CourseID string `form:"course_id" valid:"Required;"`
}

func (g GetCourseForm) GetCourseInfo(course *Course) (int, constval.ErrNo) {
	//load cache
	courseInfoCache := cache.GetGroupCache("course_info")
	if courseInfoCache == nil {
		//do not register peerpick first
		courseInfoCache = cache.NewGroupCache("course_info", maxCourseInfoCacheBytes,
			cache.GetterFunc(CourseInfoGetter))
	}
	val, err := courseInfoCache.Get(g.CourseID, cache.DefaultOption)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id": g.CourseID,
			"err":       err,
		}).Errorln("get course info cache error")
	}
	if len(val.ByteSlice()) == 0 {
		logger.GetInstance().WithField("course_id", g.CourseID).Infoln("course not exist")
		return http.StatusBadRequest, constval.CourseNotExist
	}

	//unmarshal
	err = json.Unmarshal(val.ByteSlice(), course)
	if err != nil {
		courseInfoCache.Del(g.CourseID)
		logger.GetInstance().WithField("err", err).Errorln("json unmarshal course info error")
		return http.StatusInternalServerError, constval.UnknownError
	}

	return http.StatusOK, constval.OK
}

//used for binding course
type BindCourseForm struct {
	CourseID  string `form:"course_id" valid:"Required"`
	TeacherID string `form:"teacher_id" valid:"Required"`
}

func (b BindCourseForm) BindCourse() (int, constval.ErrNo) {
	course := &Course{}
	err := Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("teacher_id").
			Where("course_id = ?", b.CourseID).First(course).Error
		if err != nil {
			return err
		}
		if course.TeacherID != nil {
			return nil
		}
		err = tx.Model(&Course{}).Select("teacher_id").Where("course_id = ?", b.CourseID).
			Update("teacher_id", b.TeacherID).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  b.CourseID,
			"teacher_id": b.TeacherID,
			"err":        err,
		}).Errorln("bind course error")
		return http.StatusInternalServerError, constval.UnknownError
	}

	if course.TeacherID != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  b.CourseID,
			"teacher_id": course.TeacherID,
		}).Infoln("course has been bound")
		return http.StatusOK, constval.CourseHasBound
	}
	return http.StatusOK, constval.OK
}

func (b BindCourseForm) UnBindCourse() (int, constval.ErrNo) {
	course := &Course{}
	err := Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("teacher_id").
			Where("course_id = ?", b.CourseID).First(course).Error
		if err != nil {
			return err
		}
		if course.TeacherID != nil {
			return nil
		}
		err = tx.Model(&Course{}).Select("teacher_id").Where("course_id = ?", b.CourseID).
			Update("teacher_id", 0).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  b.CourseID,
			"teacher_id": b.TeacherID,
			"err":        err,
		}).Errorln("unbind course error")
		return http.StatusInternalServerError, constval.UnknownError
	}

	if course.TeacherID != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id":  b.CourseID,
			"teacher_id": course.TeacherID,
		}).Infoln("course has been bound")
		return http.StatusOK, constval.CourseHasBound
	}
	return http.StatusOK, constval.OK
}

//used for getting teacher courses
type GetTeacherCourseForm struct {
	TeacherID uint64 `form:"teacher_id" valid:"Required"`
}

func (g GetTeacherCourseForm) GetTeacherCourses(courses *[]Course) (int, constval.ErrNo) {
	result := Db.Model(&Course{}).Where("teacher_id = ?", g.TeacherID).Find(courses)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.GetInstance().WithFields(logrus.Fields{
			"teacher_id": g.TeacherID,
			"err":        err,
		}).Errorln("query teacher course error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if result.RowsAffected == 0 {
		logger.GetInstance().WithField("teacher_id", g.TeacherID).Infoln("query teacher course empty")
		return http.StatusOK, constval.TeacherHasNoCourse
	}
	return http.StatusOK, constval.OK
}
