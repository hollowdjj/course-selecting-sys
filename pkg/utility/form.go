package utility

import (
	"unicode"

	"github.com/astaxie/beego/validation"
)

func PasswordCheck(v *validation.Validation, obj interface{}, key string) {
	password, ok := obj.(string)
	if !ok {
		return
	}
	var (
		hasLower = false
		hasUpper = false
		hasDigit = false
	)
	for _, c := range password {
		if hasLower && hasUpper && hasDigit {
			break
		}
		switch {
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		}
	}
	if hasLower && hasUpper && hasDigit {
		return
	}
	v.SetError("Password", "Wrong Password format")
}

func ModeCheck(v *validation.Validation, obj interface{}, key string) {
	mode, ok := obj.(string)
	if !ok {
		return
	}
	if mode == "usual" || mode == "proxy" {
		return
	}
	v.SetError("Mode", "Wrong mode")
}
