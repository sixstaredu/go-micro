package micro

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/sixstaredu/go-micro/micro/core/validate"
)
//  自定义字段验证规则
var validators = []validate.Validate{
	new(Phone),
	new(Character),
}
// 验证手机号码
type Phone struct{}

func (v *Phone) Tag() string {
	return "phone"
}
func (v *Phone) ValidateFn() validator.Func {

	return func(fl validator.FieldLevel) bool {
		return validateStrhandle(fl, func(data string) bool {
			reg := `^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199)\d{8}$`
			ok, _ := regexp.MatchString(reg, data)
			return ok
		})
	}
}
func (v *Phone) Error() string {
	return "手机号码填写不正确，请重新输入"
}
//  验证是否含有中文字符
type Character struct {}
func (*Character) Tag() string {
	return "character"
}
func (*Character) ValidateFn() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return validateStrhandle(fl, func(data string) bool {

			for _, r := range data {
				if unicode.Is(unicode.Han, r) {

					return false
				}
			}

			return true
		})
	}
}
func (*Character) Error() string {
	return "含有中文字符"
}
//
//  validateStrhandle
//  @Description: 对字符串类型统一处理方法
//  @param fl 验证字段，统一转化
//  @param fn 验证方法，将转化后的字段传入
//  @return bool
//
func validateStrhandle(fl validator.FieldLevel,fn func(filed string) bool) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return fn(data)
}
