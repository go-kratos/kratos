package validate

import (
	"regexp"

	"go-common/library/net/http/blademaster/binding"

	"gopkg.in/go-playground/validator.v9"
)

var (
	//RegIDcheck 检查身份证
	RegIDcheck = regexp.MustCompile(`(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X)$)`)
	//RegHTTPCheck 检查HTTP格式
	RegHTTPCheck = regexp.MustCompile(`^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`)
	//RegPhoneCheck 检查电话格式
	RegPhoneCheck = regexp.MustCompile(`1[345678]\d{9}`)
)

func init() {
	binding.Validator.RegisterValidation("idcheck", idcheck)
	binding.Validator.RegisterValidation("httpcheck", httpcheck)
	binding.Validator.RegisterValidation("phonecheck", phonecheck)
}

func idcheck(fl validator.FieldLevel) bool {
	return RegIDcheck.MatchString(fl.Field().String())
}

func httpcheck(fl validator.FieldLevel) bool {
	return RegHTTPCheck.MatchString(fl.Field().String())
}

func phonecheck(fl validator.FieldLevel) bool {
	return RegPhoneCheck.MatchString(fl.Field().String())
}
