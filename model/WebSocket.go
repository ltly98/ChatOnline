package model

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

/*
    编辑者：F
	功能：设置传入信息的载体
	说明：
*/

type WsCarrier struct {
	OperationCode int         `json:"opcode"`
	Carrier       interface{} `json:"carrier"`
}

/*
    编辑者：F
	功能：在线用户结构体
	说明：
*/

type OnlineUser struct {
	ID     int
	Socket *websocket.Conn
}

/*
    编辑者：F
	功能：设置操作码，方便识别和操作
	说明：
*/

const (
	//获取用户信息
	GetUserInfo = iota
	//获取聊天记录
	GetChatLog
	//发送聊天记录
	SendChatLog
	//更改聊天记录，主要是是否读取消息
	ChangeChatLog
	//获取好友请求
	GetAddRequest
	//发送好友请求
	SendAddRequest
	//修改好友请求，主要是同意或者拒绝
	ChangeAddRequest
	//获取好友分组
	GetGroup
	//修改好友分组
	ChangeGroup
	//清除好友分组
	EliminateGroup
)

//在线用户列表
var UserList = make(map[int]OnlineUser)

/*
    编辑者：F
	功能：根据操作码进行信息查询和回显
	说明：
*/

func HandledWebSocket(carr *WsCarrier) {
	switch carr.OperationCode {
	case GetUserInfo:
		wsGetUserInfo(carr)
	case GetChatLog:
		wsGetChatLog(carr)
	case SendChatLog:
		wsSendChatLog(carr)
	case ChangeChatLog:
		wsChangeChatLog(carr)
	case GetAddRequest:
		wsGetAddRequest(carr)
	case SendAddRequest:
		wsSendAddRequest(carr)
	case ChangeAddRequest:
		wsChangeAddRequest(carr)
	case GetGroup:
		wsGetGroup(carr)
	case ChangeGroup:
		wsChangeGroup(carr)
	case EliminateGroup:
		wsEliminateGroup(carr)
	default:
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码获取用户信息
	说明：测试通过
*/

func wsGetUserInfo(carr *WsCarrier) {
	var user User
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &user)
	user, err = GetUser(user.UserID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = GetUserInfo
	send["username"] = user.UserName
	send["nickname"] = user.NickName
	sendJson, err = json.Marshal(&send)
	if err != nil {
		return
	}
	err = UserList[user.UserID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码获取聊天记录
	说明：测试通过
*/

func wsGetChatLog(carr *WsCarrier) {
	var chatlog ChatLog
	var chatlogs []ChatInfo
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &chatlog)
	chatlogs, err = GetChatLogs(chatlog.SenderID, chatlog.ReceiverID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = GetChatLog
	send["data"] = chatlogs
	sendJson, err = json.Marshal(&send)
	if err != nil {
		return
	}
	err = UserList[chatlog.SenderID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码发送聊天记录
	说明：测试通过
*/
func wsSendChatLog(carr *WsCarrier) {
	var chatlog ChatLog
	var chatInfos, chatInfos2 []ChatInfo
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &chatlog)

	//存储记录,如果对方存在直接判断已读
	if _, ok := UserList[chatlog.ReceiverID]; ok {
		chatlog.LogRead = 1
	}
	err = CreateChatLog(&chatlog)
	if err != nil {
		return
	}

	//发给对方,首先判断对方是否在线
	if _, ok := UserList[chatlog.ReceiverID]; ok {
		chatInfos, err = GetChatLogs(chatlog.ReceiverID, chatlog.SenderID)
		if err != nil {
			return
		}
		send := make(map[string]interface{})
		send["datatype"] = SendChatLog
		send["data"] = chatInfos
		sendJson, _ = json.Marshal(&send)
		err = UserList[chatlog.ReceiverID].Socket.WriteMessage(1, sendJson)
	}

	//发给自己
	chatInfos2, err = GetChatLogs(chatlog.SenderID, chatlog.ReceiverID)
	if err != nil {
		return
	}
	send2 := make(map[string]interface{})
	send2["datatype"] = SendChatLog
	send2["data"] = chatInfos2
	sendJson, _ = json.Marshal(&send2)
	err = UserList[chatlog.SenderID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码更改阅读状态
	说明：测试通过
*/

func wsChangeChatLog(carr *WsCarrier) {
	var chatlog ChatLog
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &chatlog)
	err = EditChatLogRead(chatlog.SenderID, chatlog.ReceiverID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = ChangeChatLog
	send["data"] = "OK"
	sendJson, _ = json.Marshal(&send)
	//更新发送者的列表就行了
	UserList[chatlog.SenderID].Socket.WriteMessage(1, sendJson)
}

/*
    编辑者：F
	功能：根据ws操作码获取好友请求
	说明：测试通过
*/

func wsGetAddRequest(carr *WsCarrier) {
	var addRequest AddRequest
	var addRequests []RequestList
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &addRequest)
	addRequests, err = GetAddRequests(addRequest.ReceiverID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = GetAddRequest
	send["data"] = addRequests
	sendJson, err = json.Marshal(&send)
	err = UserList[addRequest.ReceiverID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码发送好友请求
	说明：测试通过
*/

func wsSendAddRequest(carr *WsCarrier) {
	var addRequest AddRequest
	var requestList []RequestList
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &addRequest)
	err = CreateAddRequest(&addRequest)
	if err != nil {
		return
	}
	requestList, err = GetAddRequests(addRequest.ReceiverID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = SendAddRequest
	send["data"] = requestList
	sendJson, err = json.Marshal(&send)
	//此处接收错误没有意义，发没发成功都要存储请求
	UserList[addRequest.ReceiverID].Socket.WriteMessage(1, sendJson)
}

/*
    编辑者：F
	功能：根据ws操作码更新好友请求,更改后删除
	说明：测试通过
*/

func wsChangeAddRequest(carr *WsCarrier) {
	var addRequest AddRequest
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &addRequest)
	send := make(map[string]interface{})
	send["datatype"] = ChangeAddRequest
	send["data"] = addRequest
	sendJson, err = json.Marshal(&send)
	//接不接收并不重要，只要数据库存储更新
	UserList[addRequest.ReceiverID].Socket.WriteMessage(1, sendJson)
	err = SwitchAddRequest(&addRequest)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码获取好友分组
	说明：测试通过
*/

func wsGetGroup(carr *WsCarrier) {
	var group UserGroup
	var groups []ContactList
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &group)
	groups, err = GetGroups(group.UserID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = GetGroup
	send["data"] = groups
	sendJson, err = json.Marshal(&send)
	err = UserList[group.UserID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码更改分组，修改过后返回分组
	说明：测试通过
*/

func wsChangeGroup(carr *WsCarrier) {
	var updateGroup UpdateGroup
	var groups []ContactList
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &updateGroup)
	err = EditGroups(&updateGroup)
	if err != nil {
		return
	}
	groups, err = GetGroups(updateGroup.UserID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = ChangeGroup
	send["data"] = groups
	sendJson, err = json.Marshal(&send)
	err = UserList[updateGroup.UserID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}

/*
    编辑者：F
	功能：根据ws操作码删除分组
	说明：测试通过
*/

func wsEliminateGroup(carr *WsCarrier) {
	var group UserGroup
	var groups []ContactList
	var err error
	var sendJson []byte
	data, _ := json.Marshal(carr.Carrier)
	json.Unmarshal(data, &group)
	err = DeleteGroups(&group)
	if err != nil {
		return
	}
	groups, err = GetGroups(group.UserID)
	if err != nil {
		return
	}
	send := make(map[string]interface{})
	send["datatype"] = EliminateGroup
	send["data"] = groups
	sendJson, err = json.Marshal(&send)
	err = UserList[group.UserID].Socket.WriteMessage(1, sendJson)
	if err != nil {
		return
	}
}
