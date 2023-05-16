package v1

import (
	"chat-online/model"
	"chat-online/msg"
	"chat-online/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	编辑者：F
	功能：发送验证码
	说明：内部自带检查，主要是同意或者拒绝的处理
*/

func SendCode(c *gin.Context) {
	//用于存储数据
	var data model.EmailParam
	data.EmailTo = c.PostForm("emailto")
	data.EmailTime = c.PostForm("emailtime")
	//数据验证
	ms, code := validator.Validate(&data)
	if code != msg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": ms,
		})
		return
	}
	//获取验证码，此处不是回显验证码，只是返回一个状态
	code = model.GenerateCode(data.EmailTo, data.EmailTime)
	//返回json
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": msg.GetErrMsg(code),
	})
}

/*
	编辑者：F
	功能：验证验证码
	说明：内部自带检查，主要是同意或者拒绝的处理
*/

func ReplaceEmail(c *gin.Context) {
	//用于存储数据
	var data model.ReplaceInfo
	//绑定json
	data.OldEmail = c.PostForm("oldemail")
	data.OldCode = c.PostForm("oldcode")
	data.NewEmail = c.PostForm("newemail")
	data.NewCode = c.PostForm("newcode")
	//数据验证
	ms, code := validator.Validate(&data)
	if code != msg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": ms,
		})
		return
	}
	//获取验证码，此处不是回显验证码，只是返回一个状态
	code = model.ReplaceEmail(&data)
	//返回json
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": "更换成功！",
	})
}
