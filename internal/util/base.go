package util

import (
	"regexp"
)

const (
	phone_regular = "^1\\d{10}$"
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
