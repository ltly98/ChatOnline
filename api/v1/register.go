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
	功能：注册用户
	说明：添加用户信息
*/

func Register(c *gin.Context) {
	//用于存储信息和回显
	var data model.RegisterInfo
	//获取表单数据
	data.UserName = c.PostForm("username")
	data.Password = c.PostForm("password")
	data.Password2 = c.PostForm("password2")
	data.Email = c.PostForm("email")
	data.EmailCode = c.PostForm("emailcode")
	//数据验证
	ms, code := validator.Validate(&data)
	if code != msg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": ms,
		})
		return
	}
	//接受自定义状态信息
	code = model.Register(&data)
	//返回json
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": "注册成功！",
	})
}
