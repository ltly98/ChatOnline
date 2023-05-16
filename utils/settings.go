package utils

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
	JwtKey   string

	DbType     string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	EmailFrom string
	EmailPwd  string
	EmailHost string
	EmailPort string
)

func InitSettings() {
	file, err := ini.Load("E:/Go/src/chat-online/config/config.ini")
	if err != nil {
		fmt.Printf("配置文件出错：%v\n", err)
	}

	LoadServer(file)
	LoadDataBase(file)
	LoadEmail(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
}

func LoadDataBase(file *ini.File) {
	DbType = file.Section("database").Key("DbType").MustString("mysql")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassword = file.Section("database").Key("DbPassword").MustString("123456")
	DbName = file.Section("database").Key("DbName").MustString("chat")
}

func LoadEmail(file *ini.File) {
	EmailFrom = file.Section("email").Key("EmailFrom").MustString("1028460024@qq.com")
	EmailPwd = file.Section("email").Key("EmailPwd").MustString("nxqozxmrvledbbee")
	EmailHost = file.Section("email").Key("EmailHost").MustString("smtp.qq.com")
	EmailPort = file.Section("email").Key("EmailPort").MustString("587")
}
