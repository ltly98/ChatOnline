package main

import (
	"chat-online/model"
	"chat-online/routes"
	"chat-online/utils"
)

func main() {
	//初始化设置
	utils.InitSettings()
	//初始化并连接数据库
	model.InitDb()
	//初始化路由
	routes.InitRouter()
}
