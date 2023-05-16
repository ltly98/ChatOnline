package msg

/*
	编辑者：F
	功能：自定义消息状态码
	说明：1-1000   	通用模块
		 1001-2000	用户模块
*/

const (

	//基本
	SUCCESS = 200
	ERROR   = 500

	//用户错误
	ERROR_USER_ID_WRONG   = 1001
	ERROR_USER_NAME_USED  = 1002
	ERROR_USER_NAME_WRONG = 1003
	ERROR_USER_NOT_EXIST  = 1004

	//密码错误
	ERROR_PASSWORD_WRONG   = 1005
	ERROR_PASSWORD_UNEQUAL = 1006

	//TOKEN错误
	ERROR_TOKEN_NOT_EXIST  = 1007
	ERROR_TOKEN_RUNTIME    = 1008
	ERROR_TOKEN_WRONG      = 1009
	ERROR_TOKEN_TYPE_WRONG = 1010

	//邮箱错误
	ERROR_EMAIL_USED      = 1011
	ERROR_EMAIL_NOT_EXIST = 1012
	ERROR_CODE_TIMEOUT    = 1013
)

var codemsg = map[int]string{
	SUCCESS:                "OK",
	ERROR:                  "FAIL",
	ERROR_USER_ID_WRONG:    "用户ID错误！",
	ERROR_USER_NAME_USED:   "用户名已存在！",
	ERROR_USER_NAME_WRONG:  "用户名错误！",
	ERROR_EMAIL_USED:       "邮箱已绑定！",
	ERROR_EMAIL_NOT_EXIST:  "邮箱不存在！",
	ERROR_CODE_TIMEOUT:     "验证码超时！",
	ERROR_PASSWORD_WRONG:   "密码错误！",
	ERROR_PASSWORD_UNEQUAL: "两次密码不相同！",
	ERROR_USER_NOT_EXIST:   "用户不存在！",
	ERROR_TOKEN_NOT_EXIST:  "TOKEN不存在！",
	ERROR_TOKEN_RUNTIME:    "TOKEN已过期！",
	ERROR_TOKEN_WRONG:      "TOKEN不正确！",
	ERROR_TOKEN_TYPE_WRONG: "TOKEN格式不正确！",
}

func GetErrMsg(code int) string {
	return codemsg[code]
}
