package v1

import (
	"chat-online/middleware"
	"chat-online/model"
	"chat-online/msg"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	编辑者：F
	功能：用户登陆
	说明：登陆成功则生成token，此处使用表单
*/

func Login(c *gin.Context) {
	//获取表单数据
	UserName := c.PostForm("username")
	PassWord := c.PostForm("password")
	var token string
	//接受自定义状态信息
	userid, code := model.Login(UserName, PassWord)
	//登陆成功才生成token
	if code == msg.SUCCESS {
		token, code = middleware.SetToken(UserName)
	}
	//返回json
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": msg.GetErrMsg(code),
		"userid":  userid,
		"token":   token,
	})

}
