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
	功能：忘记密码
	说明：内部自带检查，主要是同意或者拒绝的处理
*/

func Forget(c *gin.Context) {
	//用于存储信息
	var data model.ForgetInfo
	//绑定json
	data.Email = c.PostForm("email")
	data.EmailCode = c.PostForm("emailcode")
	data.Password = c.PostForm("password")
	data.Password2 = c.PostForm("password2")
	//数据验证
	ms, code := validator.Validate(&data)
	if code != msg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": ms,
		})
		return
	}
	//修改密码
	code = model.Forget(&data)
	//返回json
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": "修改成功！",
	})
}
