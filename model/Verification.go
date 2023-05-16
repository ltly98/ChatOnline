package model

import (
	"chat-online/msg"
	"chat-online/utils"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

/*
    编辑者：F
	功能：验证码gorm模型
	说明：
*/

type EmailVerification struct {
	VerificationID    int    `gorm:"primaryKey;type:int;not null;autoIncrement" json:"verificationid"`
	VerificationEmail string `gorm:"type:varchar(30);not null" json:"verificationemail" validate:"required,email"`
	VerificationCode  string `gorm:"type:varchar(10);not null" json:"verificationcode" validate:"required,min=6,max=10"`
	VerificationTime  string `gorm:"type:varchar(20);not null" json:"verificationtime" validate:"required,min=6,max=20"`
}

/*
    编辑者：F
	功能：此处主要是做数据验证
	说明：此处不对传入的验证码进行验证，后面进行手动验证
*/

type EmailParam struct {
	EmailTo   string `json:"emailto" validate:"required,email"`
	EmailCode int    `json:"emailcode"`
	EmailTime string `json:"emailtime" validate:"required,min=6,max=20"`
}

/*
    编辑者：F
	功能：检查邮箱验证码存储信息
	说明：检查对应邮箱的验证码是否存储，存在则返回成功
*/

func CheckCodeSave(EmailTo string) int {
	//用于存储信息和回显
	var emailVerification EmailVerification
	//执行sql，查询单列信息
	err := Db.Where("verification_email = ?", EmailTo).First(&emailVerification).Error
	if err != nil || emailVerification.VerificationID <= 0 {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
    编辑者：F
	功能：判断时间是否过期
	说明：查看是否存在记录，不存在记录直接返回成功，存在记录则对比，时间相差在允许范围内返回成功
*/

func JudgeTime(EmailTo string, EmailTime string, intervalTime int64) int {
	//设置时间格式
	formatTime := "2006/01/02 15:04:05"
	//用于存储信息和回显
	var emailVerification EmailVerification
	//执行sql，查询单列信息
	err := Db.Where("verification_email = ?", EmailTo).First(&emailVerification).Error
	if err != nil {
		return msg.ERROR
	}
	//存在时间则进行赋值、转换、判断
	//首先转换传入时间
	nowTime, err := time.ParseInLocation(formatTime, EmailTime, time.Local)
	if err != nil {
		return msg.ERROR
	}
	//然后转换存储的时间
	emailTime, err2 := time.ParseInLocation(formatTime, emailVerification.VerificationTime, time.Local)
	if err2 != nil {
		return msg.ERROR
	}
	//判断是否超过间隔时间
	if nowTime.Unix()-emailTime.Unix() > intervalTime {
		return msg.ERROR_CODE_TIMEOUT
	}
	return msg.SUCCESS
}

/*
    编辑者：F
	功能：保存获取验证码记录
	说明：一个邮箱始终保存一个验证码，暂时不删除记录
*/

func SaveCode(EmailTo string, EmailCode string, EmailTime string) int {
	//检查是否保存过该邮箱的验证码
	code := CheckCodeSave(EmailTo)
	//不存在记录，创建记录
	if code != msg.SUCCESS {
		//存储数据
		var emailVerification EmailVerification
		emailVerification.VerificationCode = EmailCode
		emailVerification.VerificationEmail = EmailTo
		emailVerification.VerificationTime = EmailTime
		//存入数据库
		err := Db.Create(&emailVerification).Error
		if err != nil {
			return msg.ERROR
		}
		return msg.SUCCESS
	}
	//存在记录，更新记录
	err := Db.Model(&EmailVerification{}).Where("verification_email = ?", EmailTo).Updates(EmailVerification{VerificationCode: EmailCode, VerificationTime: EmailTime}).Error
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
    编辑者：F
	功能：删除存储的验证码
	说明：验证过立即删除
*/

func DeleteCode(EmailTo string) int {
	//检查是否保存过该邮箱的验证码
	code := CheckCodeSave(EmailTo)
	//如果不存在，就算是已经删除了
	if code != msg.SUCCESS {
		return msg.SUCCESS
	}
	err := Db.Where("verification_email = ?", EmailTo).Delete(&EmailVerification{}).Error
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
    编辑者：F
	功能：获取验证码
	说明：
*/

func GenerateCode(EmailTo string, EmailTime string) int {
	//首先检查是否存储有验证码
	code := CheckCodeSave(EmailTo)
	if code == msg.SUCCESS {
		//判断是否超过生成间隔时间,此处必须超时
		code = JudgeTime(EmailTo, EmailTime, 60)
		if code != msg.ERROR_CODE_TIMEOUT {
			return msg.ERROR
		}
	}
	//设置base64
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	//设置验证信息
	auth := smtp.PlainAuth("", utils.EmailFrom, utils.EmailPwd, utils.EmailHost)
	//设置发送方
	to := []string{EmailTo}
	//随机种子
	rand.Seed(time.Now().Unix())
	//生成6位数的验证码
	num := rand.Intn(1000000)
	if num < 100000 {
		num += 100000
	}
	//保存
	code = SaveCode(EmailTo, fmt.Sprintf("%d", num), EmailTime)
	if code != msg.SUCCESS {
		return code
	}
	//设置发送信息header
	header := make(map[string]string)
	header["From"] = "1028460024@qq.com"
	header["To"] = EmailTo
	header["Subject"] = "验证码"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	//设置发送信息正文内容
	str := fmt.Sprintf("欢迎使用在线聊天系统，生成验证码：%v（验证有效期为5分钟，转发可能导致账号被盗！）", num)
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	//正文内容需要进行base64编码，否则接收信息为乱码
	message += "\r\n" + b64.EncodeToString([]byte(str))
	//发送邮件
	err := smtp.SendMail("smtp.qq.com:587", auth, "1028460024@qq.com", to, []byte(message))
	if err != nil {
		return msg.ERROR
	}
	return msg.SUCCESS
}

/*
    编辑者：F
	功能：验证验证码
	说明：利用当前时间进行验证，防止传输时间伪造
*/

func VerificationCode(EmailTo string, EmailCode int) int {
	//检查是否有记录
	code := CheckCodeSave(EmailTo)
	if code != msg.SUCCESS {
		return code
	}
	//判断时间
	EmailTime := time.Now().Format("2006/01/02 15:04:05")
	code = JudgeTime(EmailTo, EmailTime, 300)
	if code != msg.SUCCESS {
		return code
	}
	//用于存储信息和回显
	var emailVerification EmailVerification
	//执行sql，查询单列信息
	err := Db.Where("verification_email = ?", EmailTo).First(&emailVerification).Error
	if err != nil || emailVerification.VerificationID <= 0 {
		return msg.ERROR
	}
	SaveCode, err2 := strconv.Atoi(emailVerification.VerificationCode)
	if err2 != nil || SaveCode != EmailCode {
		return msg.ERROR
	}
	//验证完便删除，保证只验证一次
	return DeleteCode(EmailTo)
}
