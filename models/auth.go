package models

import "gorm.io/gorm"

//检查用户名和密码是否正确
func CheckAuth(username, password string) (bool, MemberInfo, error) {
	var member MemberInfo
	err := db.Model(&Member{}).Select("user_id", "user_type").
		Where("username = ? AND password = ? AND is_active = ?", username, password, 1).
		First(&member).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		//发生错误
		return false, MemberInfo{}, err
	}

	if member.UserID > 0 {
		//用户名存在且密码正确
		return true, member, nil
	}

	//用户名不存在
	return false, MemberInfo{}, nil
}
