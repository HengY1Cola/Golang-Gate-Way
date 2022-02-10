package middleware

import (
	"gin/public"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//----------------------- 自定义验证方法 -----------------------
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go

			//  检测Username是否为Admin
			val.RegisterValidation("isValidateUserName", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})

			// 验证密码是否符合格式
			val.RegisterValidation("isConformPwdFormat", func(fl validator.FieldLevel) bool {
				// 保证为8位以上的密码以及大写开头
				wantChangePwd := fl.Field().String()
				flag := true
				if len(wantChangePwd) < 8 {
					flag = false
				}
				if wantChangePwd[0] < 65 || wantChangePwd[0] > 90 {
					flag = false
				}
				return flag
			})

			//  验证服务名称
			val.RegisterValidation("validServiceName", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return matched
			})

			//  验证接入路径规则
			val.RegisterValidation("validRule", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
				return matched
			})

			//  验证URL重新规则
			val.RegisterValidation("validUrlRewrite", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, each := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(each, " ")) != 2 {
						return false
					}
				}
				return true
			})

			//  验证header转换规则
			val.RegisterValidation("validHeaderTransfor", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, each := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(each, " ")) != 3 {
						return false
					}
				}
				return true
			})

			// 验证IP列表
			val.RegisterValidation("validIpPortList", func(fl validator.FieldLevel) bool {
				for _, each := range strings.Split(fl.Field().String(), ",") {
					match, _ := regexp.Match(`^\S+\:\d+$`, []byte(each))
					if !match {
						return false
					}
				}
				return true
			})

			// 验证权重列表
			val.RegisterValidation("validWeightList", func(fl validator.FieldLevel) bool {
				for _, each := range strings.Split(fl.Field().String(), ",") {
					match, _ := regexp.Match(`^\d+$`, []byte(each))
					if !match {
						return false
					}
				}
				return true
			})

			// 验证IPList
			val.RegisterValidation("validIpList", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, item := range strings.Split(fl.Field().String(), ",") {
					matched, _ := regexp.Match(`\S+`, []byte(item)) //ip_addr
					if !matched {
						return false
					}
				}
				return true
			})

			//----------------------- 自定义验证器 -----------------------
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			val.RegisterTranslation("isValidateUserName", trans, func(ut ut.Translator) error {
				return ut.Add("isValidateUserName", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("isValidateUserName", fe.Field())
				return t
			})

			val.RegisterTranslation("isConformPwdFormat", trans, func(ut ut.Translator) error {
				return ut.Add("isConformPwdFormat", "{0} 不满足8位数以及大写字母开头", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("isConformPwdFormat", fe.Field())
				return t
			})

			val.RegisterTranslation("validServiceName", trans, func(ut ut.Translator) error {
				return ut.Add("validServiceName", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validServiceName", fe.Field())
				return t
			})

			val.RegisterTranslation("validRule", trans, func(ut ut.Translator) error {
				return ut.Add("validRule", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validRule", fe.Field())
				return t
			})

			val.RegisterTranslation("validUrlRewrite", trans, func(ut ut.Translator) error {
				return ut.Add("validUrlRewrite", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validUrlRewrite", fe.Field())
				return t
			})

			val.RegisterTranslation("validHeaderTransfor", trans, func(ut ut.Translator) error {
				return ut.Add("validHeaderTransfor", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validHeaderTransfor", fe.Field())
				return t
			})

			val.RegisterTranslation("validIpPortList", trans, func(ut ut.Translator) error {
				return ut.Add("validIpPortList", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validIpPortList", fe.Field())
				return t
			})

			val.RegisterTranslation("validWeightList", trans, func(ut ut.Translator) error {
				return ut.Add("validWeightList", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validWeightList", fe.Field())
				return t
			})

			val.RegisterTranslation("validIpList", trans, func(ut ut.Translator) error {
				return ut.Add("validIpList", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validIpList", fe.Field())
				return t
			})

			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}
