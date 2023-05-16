package model

import (
	"errors"
)

/*
	编辑者：F
	功能：聊天记录orm模型
	说明：
*/

type ChatLog struct {
	LogID      int    `gorm:"primaryKey;type:int;not null;autoIncrement" json:"logid"`
	SenderID   int    `gorm:"type:int;not null" json:"senderid"`
	ReceiverID int    `gorm:"type:int;not null" json:"receiverid"`
	LogTime    string `gorm:"type:varchar(20);not null" json:"logtime"`
	LogContent string `gorm:"type:varchar(200);not null" json:"logcontent"`
	LogRead    int    `gorm:"type:int;not null" json:"logread"`
}

/*
	编辑者：F
	功能：发送者和收者普通结构体
	说明：方便前端显示相关嵌套信息，此处结构体确有重复
*/

type ChatUser struct {
	ID       int    `json:"id"`
	NickName string `json:"nickname"`
	Remarks  string `json:"remarks"`
}

/*
	编辑者：F
	功能：回显聊天记录的结构体
	说明：方便前端显示相关嵌套信息
*/

type ChatInfo struct {
	ID       int      `json:"id"`
	Sender   ChatUser `json:"sender"`
	Receiver ChatUser `json:"receiver"`
	Time     string   `json:"time"`
	Content  string   `json:"content"`
	Read     int      `json:"read"`
}

/*
	编辑者：F
	功能：检查收发者是否存在
	说明：此处检查每一项会从效率上考虑，所以写了很多行可以合并的代码,根据成功与否返回自定义状态码
*/

func CheckChatLog(senderid int, receiverid int) error {
	//判断用户是否存在
	if senderid <= 0 || receiverid <= 0 || senderid == receiverid {
		return errors.New("用户ID出错")
	}
	//存储sender结构体信息
	var sender User
	// SELECT * FROM chat_logs WHERE sender_id = ?;
	Db.First(&sender, senderid)
	if sender.UserID <= 0 {
		return errors.New("用户不存在")
	}
	//存储receiver结构体信息
	var receiver User
	// SELECT * FROM chat_logs WHERE receiver_id = ?;
	Db.First(&receiver, senderid)
	if receiver.UserID <= 0 {
		return errors.New("用户不存在")
	}
	return nil
}

/*
	编辑者：F
	功能：获取单个聊天用户基本信息
	说明：此处senderid=0，则不查询备注名称,此处规避出错，出错直接返回空
*/

func GetChatUser(senderid int, userid int) ChatUser {
	//用于存储信息和回显
	var chatuser ChatUser
	//获取userid，此处忽略默认错误提示，使用自定义的错误
	user, _ := GetUser(userid)
	chatuser.ID = user.UserID
	chatuser.NickName = user.NickName
	//判断是获取sender还是receiver
	if senderid > 0 {
		//用于存储信息和回显
		var usergroup UserGroup
		//执行sql并接受错误
		err := Db.Where("user_id = ? AND friend_id = ?", senderid, userid).First(&usergroup).Error
		if err == nil {
			chatuser.Remarks = usergroup.Remarks
		}
	}
	return chatuser
}

/*
	编辑者：F
	功能：重组为指定结构体信息存储
	说明：此处直接赋值，不返回错误信息，因为都是聊天记录里已有的信息
*/

func DivideChatLog(chatlogs []ChatLog) []ChatInfo {
	//用于存储信息
	var chatinfos []ChatInfo
	//遍历传入的切片
	for _, value := range chatlogs {
		//用于存储信息
		var chatinfo ChatInfo
		chatinfo.ID = value.LogID
		chatinfo.Sender = GetChatUser(0, value.SenderID)
		chatinfo.Receiver = GetChatUser(value.SenderID, value.ReceiverID)
		chatinfo.Time = value.LogTime
		chatinfo.Content = value.LogContent
		chatinfo.Read = value.LogRead
		//添加进存储的新切片
		chatinfos = append(chatinfos, chatinfo)
	}
	return chatinfos
}

/*
	编辑者：F
	功能：创建聊天记录
	说明：根据成功与否返回自定义状态码
*/

func CreateChatLog(chatlog *ChatLog) error {
	//检查用户分组，获取自定义状态信息
	err := CheckChatLog(chatlog.SenderID, chatlog.ReceiverID)
	if err != nil {
		return err
	}
	//执行sql并接受错误
	err = Db.Create(&chatlog).Error
	if err != nil {
		return err
	}
	return err
}

/*
    编辑者：F
	功能：获取聊天记录
	说明：此处获取过后，处理一个数组返回
*/

func GetChatLogs(senderid int, receiverid int) ([]ChatInfo, error) {
	//检查用户分组，获取自定义状态信息
	err := CheckChatLog(senderid, receiverid)
	if err != nil {
		return nil, err
	}
	//用于存储信息和回显
	var chatlogs []ChatLog
	//执行sql并接受错误
	err = Db.Where("sender_id in (?,?) AND receiver_id in (?,?)", senderid, receiverid, senderid, receiverid).Order("log_time").Find(&chatlogs).Error
	if err != nil {
		return nil, err
	}
	//处理用户分组
	chatinfos := DivideChatLog(chatlogs)
	return chatinfos, nil
}

/*
	编辑者：F
	功能：主要用于接收者是否阅读消息
	说明：根据成功与否返回自定义状态码,未读默认传入0，已读传入1
*/

func EditChatLogRead(senderid int, receiverid int) error {
	//检查用户分组，获取自定义状态信息
	err := CheckChatLog(senderid, receiverid)
	if err != nil {
		return err
	}
	//执行sql并接受错误
	err = Db.Model(&ChatLog{}).Where("sender_id = ? AND receiver_id = ? AND log_read = 0", senderid, receiverid).Update("log_read", 1).Error
	if err != nil {
		return err
	}
	return nil
}
