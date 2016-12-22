package util

import (
	"regexp"
)

const (
	phone_regular = "^1(3[0-9]|4[57]|5[0-35-9]|7[0135678]|8[0-9])\\d{8}$"
	email_regular = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
)

func ValidateMob(mobileNum string) bool {
	return validate(mobileNum, phone_regular)
}

func ValidateEmail(email string) bool {
	return validate(email, email_regular)
}

func validate(val, regular string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(val)
}
