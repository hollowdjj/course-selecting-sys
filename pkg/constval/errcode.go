package constval

type ErrNo int

const (
	OK ErrNo = iota

	//for login
	UserExisted
	UserDeleted
	UserNotExist
	WrongPassword
	LoginRequired
	PermDenied

	UserDeletedOrNotExist

	//for course
	CourseExisted
	CourseNotExist
	CourseNotAvailable
	CourseHasBound
	CourseNotBind
	UnBindError
	TeacherHasNoCourse

	//for course selecting
	StudentNotExist
	StudentHasNoCourse
	StudentHasCourse

	ParamInvalid
	UnknownError
)

var msg = map[ErrNo]string{
	OK: "ok",

	UserExisted:   "Username已存在",
	UserDeleted:   "用户已删除",
	UserNotExist:  "用户不存在",
	WrongPassword: "密码错误",
	LoginRequired: "用户未登录",
	PermDenied:    "没有操作权限",

	UserDeletedOrNotExist: "用户不存在或已删除",

	CourseExisted:      "课程已存在",
	CourseNotExist:     "课程不存在",
	CourseNotAvailable: "课程已满",
	CourseHasBound:     "课程已绑定过",
	CourseNotBind:      "课程未绑定过",
	UnBindError:        "课程绑定的不是该老师",
	TeacherHasNoCourse: "老师没有该课程",

	StudentNotExist:    "学生不存在",
	StudentHasNoCourse: "学生没有选择任何课程",
	StudentHasCourse:   "学生有课程",

	ParamInvalid: "参数不合法",
	UnknownError: "未知错误",
}

//ger error message according to err code
func GetErrCodeMsg(code ErrNo) string {
	return msg[code]
}
