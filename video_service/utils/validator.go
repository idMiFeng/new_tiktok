package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

const DefaultLocale = "zh_Hans_CN"

var trans ut.Translator

func GetTrans() ut.Translator {
	return trans
}

func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhLocale, zhHansCNLocale := zh.New(), zh_Hans_CN.New()
		uni := ut.New(zhHansCNLocale, zhHansCNLocale, zhLocale)

		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		_ = zhTrans.RegisterDefaultTranslations(v, trans)

		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			count := 2
			name := strings.SplitN(field.Tag.Get("label"), ",", count)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		return nil
	} else {
		return errors.New("get binding.Validator.Engine() failed")
	}
}
