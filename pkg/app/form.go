package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/hollowdjj/course-selecting-sys/pkg/e"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"net/http"
)

type CustomFunc map[string]func(v *validation.Validation, obj interface{}, key string)

//表单验证。返回http状态码以及参数检验的错误代码
func BindAndValid(c *gin.Context, form interface{}, json bool) (int, e.ErrNo) {
	var err error
	if json {
		err = c.BindJSON(form)
	} else {
		err = c.BindQuery(form)
	}

	if err != nil {
		return http.StatusBadRequest, e.ParamInvalid
	}

	valid := validation.Validation{}
	right, err := valid.Valid(form)
	if err != nil {
		logging.Error(err)
		return http.StatusInternalServerError, e.UnknownError
	}
	if !right {
		MakeErrors(valid.Errors)
		return http.StatusBadRequest, e.ParamInvalid
	}

	return http.StatusOK, e.OK
}

//支持自定义表单验证函数
func BindAndValidCustom(c *gin.Context, form interface{}, funcs CustomFunc, json bool) (int, e.ErrNo) {
	var err error
	if json {
		err = c.BindJSON(form)
	} else {
		err = c.BindQuery(form)
	}

	if err != nil {
		return http.StatusBadRequest, e.ParamInvalid
	}

	//设置自定义的表单验证函数
	valid := validation.Validation{}
	for k, v := range funcs {
		if err = validation.AddCustomFunc(k, v); err != nil {
			return http.StatusInternalServerError, e.UnknownError
		}
	}

	//表单验证
	check, err := valid.Valid(form)
	if err != nil {
		logging.Error(err)
		return http.StatusInternalServerError, e.UnknownError
	}
	if !check {
		MakeErrors(valid.Errors)
		return http.StatusBadRequest, e.ParamInvalid
	}

	return http.StatusOK, e.OK
}
