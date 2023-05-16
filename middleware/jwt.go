package middleware

import (
	"chat-online/msg"
	"chat-online/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

/*
	编辑者：F,Z
	功能：指定加密密钥
	说明：用作token生成
*/

var JwtKey = []byte(utils.JwtKey)

/*
	编辑者：F,Z
	功能：Claims结构体
	说明：用作token生成
*/

type UserClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

/*
	编辑者：F,Z
	功能：生成token
	说明：Unix()添加时间戳
*/

func SetToken(username string) (string, int) {
	//设置有效时间
	expireTime := time.Now().Add(10 * time.Hour)
	SetClaims := UserClaims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			//设置过期时间，这里添加了时间戳
			ExpiresAt: expireTime.Unix(),
			//指定token发行人
			Issuer: "chat-online",
		},
	}
	//使用HS256才能加盐解析
	reqClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	//内部生成签名字符串
	token, err := reqClaims.SignedString(JwtKey)
	if err != nil {
		return "", msg.ERROR
	}
	return token, msg.SUCCESS
}

/*
	编辑者：F,Z
	功能：验证token
	说明：参考官方文档的固定写法，不要轻易修改
*/

func CheckToken(token string) (*UserClaims, int) {
	setToken, _ := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if claims, ok := setToken.Claims.(*UserClaims); ok && setToken.Valid {
		return claims, msg.SUCCESS
	} else {
		return nil, msg.ERROR
	}
}

/*
	编辑者：F,Z
	功能：jwt中间件
	说明：调用 Abort 以确保不调用此请求的剩余处理程序
*/

func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		//保证两种方式都可以
		if tokenHeader == "" {
			tokenHeader = c.Query("token")
		}
		code := msg.SUCCESS
		//如果token头部提示为空
		if tokenHeader == "" {
			code = msg.ERROR_TOKEN_NOT_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": msg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//提交格式为：Bearer token
		checkToken := strings.Split(tokenHeader, " ")
		//如果token为空
		if len(checkToken) == 0 {
			code = msg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": msg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//如果token内容不正确
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code = msg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": msg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//检查token
		key, tCode := CheckToken(checkToken[1])
		if tCode == msg.ERROR {
			code = msg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": msg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//token是否过期
		if time.Now().Unix() > key.ExpiresAt {
			code = msg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": msg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		c.Set("username", key.UserName)
		c.Next()
	}
}
