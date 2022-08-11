package models

//user type
type UserType int

const (
	Admin UserType = iota
	Student
	Teacher
)

//table user
type User struct {
	UserID   uint64 `gorm:"primaryKey" json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	UserType int    `json:"user_type"`
	IsActive int    `json:"is_active" gorm:"default:1"`
}

//query user info
type UserInfo struct {
	UserID   uint64 `json:"user_id"`
	UserType int    `json:"user_type"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

//table course
type Course struct {
	CourseID    uint64  `gorm:"primaryKey" json:"course_id"`
	CourseName  string  `json:"course_name"`
	Cap         uint    `json:"cap"`
	RemainCap   uint    `json:"remain_cap"`
	TeacherID   *uint64 `json:"teacher_id"`
	TeacherName *string `json:"teacher_name"`
}

type StudentCourse struct {
	StudentID uint64 `json:"student_id"`
	CourseID  uint64 `json:"course_id"`
}
