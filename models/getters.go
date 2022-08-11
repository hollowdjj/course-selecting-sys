package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//cache Getter of user
func UserInfoGetter(username string) ([]byte, error) {
	//query
	user := &User{}
	err := Db.Where("username = ?", username).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.GetInstance().WithFields(logrus.Fields{
			"username": username,
			"err":      err,
		}).Errorln("query user info by username error")
		return nil, err
	}
	if user.UserID == 0 {
		return nil, nil
	}
	//marshal
	data, err := json.Marshal(user)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"userinfo": *user,
			"err":      err,
		}).Errorln("json marshal user info error")
		return nil, err
	}
	return data, nil
}

func UserInfoGetterByUserID(id string) ([]byte, error) {
	//query
	userInfo := &UserInfo{}
	err := Db.Model(&User{}).Where("userid = ?", id).First(&userInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.GetInstance().WithFields(logrus.Fields{
			"userid": id,
			"err":    err,
		}).Errorln("query user info by id error")
		return nil, err
	}
	if userInfo.UserID == 0 {
		return nil, nil
	}
	//marshal
	data, err := json.Marshal(userInfo)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"userinfo": *userInfo,
			"err":      err,
		}).Errorln("json marshal user info error")
		return nil, err
	}
	return data, nil
}

//cache Getter of course
func CourseInfoGetter(courseID string) ([]byte, error) {
	course := &Course{}
	result := Db.Where("course_id = ?", courseID).First(course)
	if result.RowsAffected == 0 {
		logger.GetInstance().WithField("course_id", courseID).Infoln("course not exist")
		return nil, nil
	}
	if err := result.Error; err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id": courseID,
			"err":       err,
		}).Errorln("query course info error")
		return nil, err
	}

	data, err := json.Marshal(course)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"course": course,
			"err":    err,
		}).Errorln("json marshal course info error")
		return nil, err
	}

	return data, nil
}

//cache Getter for course remain cap
func CourseRemainCapGetter(courseID string) ([]byte, error) {
	course := &Course{}
	result := Db.Select("remain_cap").Where("course_id = ?", courseID).First(course)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.GetInstance().WithFields(logrus.Fields{
			"course_id": courseID,
			"err":       err,
		}).Errorln("query course remain cap error")
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return []byte(strconv.Itoa(int(course.RemainCap))), nil
}

//cache Getter for student course. The output depends on key:
//1. key is userid_courseid, return "1" if user has the course
//2. key is userid, return all student courses format as courseid1_courseid2....
func StudentCourseGetter(key string) ([]byte, error) {
	queryCourse, err := regexp.MatchString(`^[0-9]*$`, key)
	if err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Errorln("check key pattern error")
		return nil, err
	}
	//get all student courses
	if queryCourse {
		studentCourses := []StudentCourse{}
		err := Db.Model(&Course{}).Select("course_id").Where("student_id = ?", key).Find(&studentCourses).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.GetInstance().WithFields(logrus.Fields{
				"student_id": key,
				"err":        err,
			}).Errorln("query student courses error")
			return nil, err
		}
		ret := ""
		for i := 0; i < len(studentCourses)-1; i++ {
			ret += fmt.Sprintf("%d", studentCourses[i].CourseID) + "_"
		}
		ret += fmt.Sprintf("%d", studentCourses[len(studentCourses)-1].CourseID)
		logger.GetInstance().WithFields(logrus.Fields{
			"student_id": key,
			"courses:":   ret,
		}).Info("query student courses succ")
		return []byte(ret), nil
	}

	valid, err := regexp.MatchString(`^[0-9]*_[0-9]*$`, key)
	if !valid || err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Errorln("wrong key format")
		return nil, err
	}

	strs := strings.Split(key, "_")
	studentCourse := &StudentCourse{}
	result := Db.Where("student_id = ? AND course_id = ?", strs[0], strs[1]).First(studentCourse)
	if result.RowsAffected == 0 {
		//suggest that student does not have this course
		return nil, nil
	}
	if err := result.Error; err != nil {
		logger.GetInstance().WithFields(logrus.Fields{
			"student_id": strs[0],
			"err":        err,
		}).Errorln("get student course error")
		return nil, err
	}

	return []byte("1"), nil
}
