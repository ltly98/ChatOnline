package model

import (
	"errors"
)

/*
	编辑者：F
	功能：添加好友请求的orm模型
	说明：
*/

type AddRequest struct {
	RequestID   int    `gorm:"primaryKey;type:int;not null;autoIncrement" json:"requestid"`
	SenderID    int    `gorm:"type:int;not null" json:"senderid"`
	ReceiverID  int    `gorm:"type:int;not null" json:"receiverid"`
	RequestTime string `gorm:"type:varchar(20);not null" json:"requesttime"`
	IsAccept    int    `gorm:"type:int;not null" json:"isaccept"`
}

/*
	编辑者：F
	功能：返回json数据使用的结构体
	说明：
*/

type RequestList struct {
	RequestID   int    `json:"requestid"`
	SenderID    int    `json:"senderid"`
	SenderName  string `json:"sendername"`
	ReceiverID  int    `json:"receiverid"`
	RequestTime string `json:"requesttime"`
	IsAccept    int    `json:"isaccept"`
}

/*
    编辑者：F
	功能：检查对应的ID是否存在
	说明：
*/

func CheckUserID(senderid int, receiverid int) error {
	//不能添加自己为好友
	if senderid == receiverid {
		return errors.New("不能添加自己为好友")
	}
	var Sender, Receiver User
	//执行sql，查询发送者是否存在
	Db.Select("user_id").Where("user_id = ?", senderid).First(&Sender)
	if Sender.UserID <= 0 {
		return errors.New("发送者不存在")
	}
	//执行sql，查询接受者是否存在
	Db.Select("user_id").Where("user_id = ?", receiverid).First(&Receiver)
	if Sender.UserID <= 0 {
		return errors.New("接收者不存在")
	}
	return nil
}

/*
    编辑者：F
	功能：检查添加请求是否存在
	说明：此处不再检查用户是否存在，所以注意提前检查id
*/

func CheckAddRequest(senderid int, receiverid int) error {
	//用于存储信息
	var addRequest AddRequest
	//执行sql，查询接受者是否存在
	Db.Select("request_id").Where("sender_id = ? AND receiver_id = ?", senderid, receiverid).First(&addRequest)
	if addRequest.RequestID > 0 {
		return errors.New("已经发送过请求")
	}
	return nil
}

/*
    编辑者：F
	功能：sender和receiver各自添加对方至好友列表,默认分组为新好友
	说明：接收两个id即可
*/

func CreateGroupForRequest(senderid int, receiverid int) error {
	//检查sender新好友分组是否存在
	err := CheckGroupName(senderid, "新好友")
	//不存在则创建分组
	if err == nil {
		//该方式为创建新分组
		err = CreateGroup(&UserGroup{
			GroupName: "新好友",
			UserID:    senderid,
			FriendID:  0,
		})
		if err != nil {
			return err
		}
	}
	//检查receiver新好友分组是否存在
	err = CheckGroupName(receiverid, "新好友")
	//不存在则创建分组
	if err == nil {
		//该方式为创建新分组
		err = CreateGroup(&UserGroup{
			GroupName: "新好友",
			UserID:    receiverid,
			FriendID:  0,
		})
		if err != nil {
			return err
		}
	}
	err = CreateGroup(&UserGroup{
		GroupName: "新好友",
		UserID:    senderid,
		FriendID:  receiverid,
		Remarks:   "",
	})
	if err != nil {
		return err
	}
	err = CreateGroup(&UserGroup{
		GroupName: "新好友",
		UserID:    receiverid,
		FriendID:  senderid,
		Remarks:   "",
	})
	if err != nil {
		return err
	}
	return nil
}

/*
    编辑者：F
	功能：获取好友请求列表
	说明：将添加好友请求处理后，返回新对象的请求列表
*/

func GetRequestList(addRequests []AddRequest) ([]RequestList, error) {
	//存储结果
	var requestList []RequestList
	//存储单个结果
	var simpleRequest RequestList
	//遍历处理
	for _, value := range addRequests {
		user, err := GetUser(value.SenderID)
		if err != nil {
			return nil, err
		}
		simpleRequest.RequestID = value.RequestID
		simpleRequest.SenderID = value.SenderID
		simpleRequest.SenderName = user.NickName
		simpleRequest.ReceiverID = value.ReceiverID
		simpleRequest.RequestTime = value.RequestTime
		simpleRequest.IsAccept = value.IsAccept
		requestList = append(requestList, simpleRequest)
	}
	return requestList, nil
}

/*
    编辑者：F
	功能：获取添加好友记录
	说明：接收自己的id即可
*/

func GetAddRequests(receiverid int) ([]RequestList, error) {
	var addRequests []AddRequest
	if receiverid <= 0 {
		return nil, errors.New("传入ID不正确")
	}
	//执行sql语句并返回错误信息
	err := Db.Where("receiver_id = ?", receiverid).Find(&addRequests).Error
	if err != nil {
		return nil, err
	}
	requestList, err2 := GetRequestList(addRequests)
	if err2 != nil {
		return nil, err2
	}
	return requestList, nil
}

/*
    编辑者：F
	功能：创建添加好友记录
	说明：接收两个id和时间即可
*/

func CreateAddRequest(add *AddRequest) error {
	//检查用户的ID
	err := CheckUserID(add.SenderID, add.ReceiverID)
	if err != nil {
		return err
	}
	//检查是否已经发送过请求
	err = CheckAddRequest(add.SenderID, add.ReceiverID)
	if err != nil {
		return err
	}
	//创建结构体存储数据
	var addRequest AddRequest
	addRequest.SenderID = add.SenderID
	addRequest.ReceiverID = add.ReceiverID
	addRequest.RequestTime = add.RequestTime
	addRequest.IsAccept = -1
	//存储结构体数据
	err = Db.Create(&addRequest).Error
	if err != nil {
		return err
	}
	return nil
}

/*
    编辑者：F
	功能：更新添加好友的请求
	说明：接收两个id和isaccept即可，isaccept默认为-1，同意为1，拒绝为0,处理完直接删除记录
*/

func SwitchAddRequest(add *AddRequest) error {
	//检查用户的ID
	err := CheckUserID(add.SenderID, add.ReceiverID)
	if err != nil {
		return err
	}
	err = Db.Where("request_id = ? ", add.RequestID).Delete(&AddRequest{}).Error
	if err != nil {
		return err
	}
	//同意则添加至新好友分组
	if add.IsAccept == 1 {
		err = CreateGroupForRequest(add.SenderID, add.ReceiverID)
		if err != nil {
			return err
		}
	}
	return nil
}
