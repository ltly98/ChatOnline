package validator

import (
	"chat-online/msg"
	"fmt"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

func Validate(data interface{}) (string, int) {
	validate := validator.New()
	//此处主要用作将提示信息转换为中文，如果不需要可以根据相关进行删除
	uni := unTrans.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN")
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err:", err)
	}
	//判定结构体
	err2 := validate.Struct(data)
	if err2 != nil {
		for _, v := range err2.(validator.ValidationErrors) {
			return v.Translate(trans), msg.ERROR
		}
	}
	return "", msg.SUCCESS
}
