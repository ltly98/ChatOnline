package model

import (
	"chat-online/msg"
	"errors"
	"strconv"
)

/*
	编辑者：F,W
	功能：用户表orm模型结构体
	说明：此处添加数据验证
*/

type User struct {
	UserID   int    `gorm:"primaryKey;type:int;not null;autoIncrement" json:"userid"`
	UserName string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=6,max=18"`
	PassWord string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=18"`
	Email    string `gorm:"type:varchar(30);not null" json:"email" validate:"required,email"`
	NickName string `gorm:"type:varchar(10);not null" json:"nickname"`
}

/*
	编辑者：F
	功能：注册解析json使用结构体
	说明：
*/

type RegisterInfo struct {
	UserName  string `json:"username" validate:"required,min=6,max=18"`
	Password  string `json:"password" validate:"required,min=6,max=18"`
	Password2 string `json:"password2" validate:"required,min=6,max=18"`
	Email     string `json:"email" validate:"required,email"`
	EmailCode string `json:"emailcode" validate:"required,len=6"`
}

/*
	编辑者：F
	功能：忘记密码解析json使用结构体
	说明：
*/

type ForgetInfo struct {
	Email     string `json:"email" validate:"required,email"`
	EmailCode string `json:"emailcode" validate:"required,len=6"`
	Password  string `json:"password" validate:"required,min=6,max=18"`
	Password2 string `json:"password2" validate:"required,min=6,max=18"`
}

/*
	编辑者：F
	功能：更换邮箱解析json使用结构体
	说明：
*/

type ReplaceInfo struct {
	OldEmail string `json:"oldemail" validate:"required,email"`
	OldCode  string `json:"oldcode" validate:"required,len=6"`
	NewEmail string `json:"newemail" validate:"required,email"`
	NewCode  string `json:"newcode" validate:"required,len=6"`
}

/*
	编辑者：F,W
	功能：使用gorm检查用户是否存在
	说明：此处是注册用户使用,用户必须不存在
*/

func CheckUserByUserName(name string) int {
	//用于存储信息和回显
	var users User
	//执行sql，查询单列信息
	Db.Select("user_id").Where("user_name = ?", name).First(&users)
	if users.UserID > 0 {
		return msg.ERROR_USER_NAME_USED
	}
	return msg.SUCCESS
}

/*
	编辑者：F,W
	功能：使用gorm检查用户是否存在
	说明：忘记密码使用时，应该返回成功，更换邮箱使用时，应该返回错误用作判断
*/

func CheckUserByEmail(email string) int {
	//用于存储信息和回显
	var users User
	//执行sql，查询单列信息
	Db.Select("user_id").Where("email = ?", email).First(&users)
	if users.UserID <= 0 {
		return msg.ERROR_EMAIL_NOT_EXIST
	}
	return msg.SUCCESS
}

/*
    编辑者：F
	功能：检查邮箱
	说明：检查绑定邮箱是否注册
*/

func CheckEmail(EmailTo string) int {
	//用于存储信息和回显
	var users User
	//执行sql，查询单列信息
	Db.Select("user_id").Where("email = ?", EmailTo).First(&users)
	if users.UserID > 0 {
		return msg.ERROR_EMAIL_USED
	}
	return msg.SUCCESS
}

/*
	编辑者：F
	功能：检查密码是否相同
	说明：直接对比返回结果
*/

func CheckPasswordEqual(password string, password2 string) int {
	if password != password2 {
		return msg.ERROR_PASSWORD_UNEQUAL
	}
	return msg.SUCCESS
}

/*
	编辑者：F,W
	功能：使用gorm获取指定用户信息
	说明：通过gorm查询指定id的用户信息，成功返回user，失败返回nil，都会将自定义消息状态码返回
*/

func GetUser(ID int) (User, error) {
	//用于存储信息和回显
	var user User
	//判断用户id格式
	if ID <= 0 {
		return user, errors.New("用户id不正确！")
	}
	// 根据主键检索
	// SELECT * FROM users WHERE id = ID;
	err := Db.First(&user, ID).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

/*
	编辑者：F,W
	功能：使用gorm修改用户昵称
	说明：通过gorm查询指定id的用户信息，然后进行修改，并返回自定义状态码
*/

func EditUserNickName(data *User) int {
	//获取自定义状态码
	code := CheckUserByUserName(data.UserName)
	if code != msg.SUCCESS {
		return code
	}
	//执行sql，获取错误信息
	err := Db.Model(&User{}).Where("user_id = ?", data.UserID).Update("nick_name", data.NickName).Error
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
	编辑者：F
	功能：检查登陆
	说明：检查用户名和密码是否正确,此处返回userid和自定义状态码
*/

func Login(username string, password string) (int, int) {
	//用于存储信息和回显
	var user User
	//执行sql语句，存储于user
	Db.Where("user_name = ?", username).First(&user)
	//判断用户id格式
	if user.UserID <= 0 {
		return 0, msg.ERROR_USER_NOT_EXIST
	}
	//判断密码是否正确
	if ScryptPw(password) != user.PassWord {
		return 0, msg.ERROR_PASSWORD_WRONG
	}
	return user.UserID, msg.SUCCESS
}

/*
	编辑者：F,W
	功能：注册用户
	说明：
*/

func Register(data *RegisterInfo) int {
	//检查邮箱是否绑定
	code := CheckEmail(data.Email)
	if code != msg.SUCCESS {
		return code
	}
	//检查两次密码是否相同
	code = CheckPasswordEqual(data.Password, data.Password2)
	if code != msg.SUCCESS {
		return code
	}
	//检查用户信息，获取自定义状态码
	code = CheckUserByUserName(data.UserName)
	if code != msg.SUCCESS {
		return code
	}
	//格式转换
	emailCode, err := strconv.Atoi(data.EmailCode)
	if err != nil {
		return msg.ERROR
	}
	//验证验证码
	code = VerificationCode(data.Email, emailCode)
	if code != msg.SUCCESS {
		return code
	}
	//此处默认昵称为账户名
	user := User{
		UserName: data.UserName,
		PassWord: ScryptPw(data.Password),
		Email:    data.Email,
		NickName: data.UserName,
	}
	//执行sql，获取错误信息
	err = Db.Create(&user).Error
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
	编辑者：F,W
	功能：使用gorm注册用户
	说明：
*/

func Forget(data *ForgetInfo) int {
	//检查邮箱是否绑定
	code := CheckEmail(data.Email)
	if code != msg.ERROR_EMAIL_USED {
		return msg.ERROR_USER_NOT_EXIST
	}
	//检查密码是否一致
	code = CheckPasswordEqual(data.Password, data.Password2)
	if code != msg.SUCCESS {
		return code
	}
	//检查用户信息，获取自定义状态码
	code = CheckUserByEmail(data.Email)
	if code != msg.SUCCESS {
		return code
	}
	//格式转换
	emailCode, err := strconv.Atoi(data.EmailCode)
	if err != nil {
		return msg.ERROR
	}
	//验证验证码
	code = VerificationCode(data.Email, emailCode)
	if code != msg.SUCCESS {
		return code
	}
	//执行sql，获取错误信息
	err = Db.Model(&User{}).Where("email = ?", data.Email).Update("pass_word", ScryptPw(data.Password)).Error
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
	编辑者：F,W
	功能：使用gorm修改用户绑定邮箱
	说明：新旧邮箱都要检查
*/

func ReplaceEmail(data *ReplaceInfo) int {
	//检查旧邮箱，必须存在
	code := CheckUserByEmail(data.OldEmail)
	if code != msg.SUCCESS {
		return code
	}
	//检查新邮箱，必须不存在存在
	code = CheckUserByEmail(data.NewEmail)
	if code != msg.ERROR_EMAIL_NOT_EXIST {
		return msg.ERROR_EMAIL_USED
	}
	//格式转换
	oldCode, err := strconv.Atoi(data.OldCode)
	if err != nil {
		return msg.ERROR
	}
	//验证验证码
	code = VerificationCode(data.OldEmail, oldCode)
	if code != msg.SUCCESS {
		return code
	}
	newCode, err2 := strconv.Atoi(data.NewCode)
	if err2 != nil {
		return msg.ERROR
	}
	//验证验证码
	code = VerificationCode(data.NewEmail, newCode)
	if code != msg.SUCCESS {
		return code
	}
	//执行sql，获取错误信息
	err = Db.Model(&User{}).Where("email = ?", data.OldEmail).Update("email", data.NewEmail).Error
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}
