package app

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/pkg/constval"
	"github.com/hollowdjj/course-selecting-sys/pkg/logger"
	"github.com/sirupsen/logrus"
)

type CustomFunc map[string]func(v *validation.Validation, obj interface{}, key string)

//bind and validate form. json indicates the form is in json format
func BindAndValid(c *gin.Context, form interface{}, json bool) (int, constval.ErrNo) {
	var err error
	if json {
		err = c.BindJSON(form)
	} else {
		err = c.BindQuery(form)
	}

	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("bind request form error")
		return http.StatusBadRequest, constval.ParamInvalid
	}

	valid := validation.Validation{}
	right, err := valid.Valid(form)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("validate request form error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if !right {
		MakeErrors(valid.Errors)
		return http.StatusBadRequest, constval.ParamInvalid
	}

	return http.StatusOK, constval.OK
}

//bind and validate form with custom function
func BindAndValidCustom(c *gin.Context, form interface{}, funcs CustomFunc, json bool) (int, constval.ErrNo) {
	var err error
	if json {
		err = c.BindJSON(form)
	} else {
		err = c.BindQuery(form)
	}

	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("bind request form error")
		return http.StatusBadRequest, constval.ParamInvalid
	}

	//setup custom validate function
	valid := validation.Validation{}
	for k, v := range funcs {
		if err = validation.AddCustomFunc(k, v); err != nil {
			logger.GetInstance().WithField("err", err).Errorln("set custom validate func error")
			return http.StatusInternalServerError, constval.UnknownError
		}
	}

	//validate form
	check, err := valid.Valid(form)
	if err != nil {
		logger.GetInstance().WithField("err", err).Errorln("validate request form error")
		return http.StatusInternalServerError, constval.UnknownError
	}
	if !check {
		MakeErrors(valid.Errors)
		return http.StatusBadRequest, constval.ParamInvalid
	}

	return http.StatusOK, constval.OK
}

//print form validate error message
func MakeErrors(errors []*validation.Error) {
	fields := make(map[string]interface{})
	for _, err := range errors {
		fields[err.Key] = err.Message
	}
	logger.GetInstance().WithFields(logrus.Fields(fields)).Errorln("wrong form")
}
