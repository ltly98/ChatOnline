package model

import (
	"errors"
	"log"
)

/*
	编辑者：F,W
	功能：用户分组表orm模型结构体
	说明：在创建用户的json请求中，将json的groupid设置为“”，能够正常在数据库中进行创建
		 此处做一个约定，新创建的分组有一条作为标记的记录，FriendID指定为0
*/

type UserGroup struct {
	GroupID   int    `gorm:"primaryKey;type:int;not null;autoIncrement" json:"groupid"`
	GroupName string `gorm:"type:varchar(20);not null" json:"groupname"`
	UserID    int    `gorm:"type:int;not null" json:"userid"`
	FriendID  int    `gorm:"type:int;not null" json:"friendid"`
	Remarks   string `gorm:"type:varchar(20)" json:"remarks"`
}

/*
	编辑者：F
	功能：好友普通结构体
	说明：方便前端显示相关嵌套信息
*/

type ContactList struct {
	GroupName string `json:"groupname"`
	ID        int    `json:"id"`
	NickName  string `json:"nickname"`
	Remarks   string `json:"remarks"`
}

/*
	编辑者：F
	功能：方便存储修改分组内容的信息
	说明：因为在修改分组名称的时候会存在使用旧分组名称查找来修改新分组名称，现有结构体无法满足需求
*/

type UpdateGroup struct {
	UserID       int    `json:"userid"`
	FriendID     int    `json:friendid`
	Remarks      string `json:remarks`
	OldGroupName string `json:"oldgroupname"`
	NewGroupName string `json:"newgroupname"`
}

/*
	编辑者：F
	功能：处理返回联系人的数据
	说明：之前编写的返回数据虽然更加合理，但是前端接收数据无法解析，所以更改为该方法
*/

func GetContactList(usergroups []UserGroup) []ContactList {
	//存储结果
	var contactList []ContactList
	//存储单个联系人
	var simpleContact ContactList
	//遍历处理
	for _, value := range usergroups {
		//尝试获取好友信息，获取成功就进行处理
		user, err := GetUser(value.FriendID)
		if err == nil {
			//存储好友信息
			simpleContact.GroupName = value.GroupName
			simpleContact.ID = value.FriendID
			simpleContact.NickName = user.NickName
			simpleContact.Remarks = value.Remarks
			//加入已有分组
			contactList = append(contactList, simpleContact)
		} else {
			simpleContact.GroupName = value.GroupName
			simpleContact.ID = 0
			simpleContact.NickName = ""
			simpleContact.Remarks = ""
			//加入已有分组
			contactList = append(contactList, simpleContact)
		}
	}
	return contactList
}

/*
	编辑者：F,W
	功能：使用gorm检查分组名称
	说明：不能添加自己为好友，也就是分组内userid不能和friendid相同，不允许相同用户拥有同名分组
*/

func CheckGroupName(userid int, groupname string) error {
	//用于存储信息和回显
	var group UserGroup
	//执行sql语句并存储在group中
	Db.Select("group_id").Where("user_id = ? AND group_name = ?", userid, groupname).First(&group)
	if group.GroupID > 0 {
		return errors.New("分组名已存在")
	}
	return nil
}

/*
	编辑者：F,W
	功能：使用gorm检查分组名称
	说明：规定好友只能存在于一个分组，并且不能加自己为好友
*/

func CheckGroupFriendID(userid int, friendid int) error {
	//不能添加自己为好友
	if userid == friendid {
		return errors.New("不能添加自己为好友")
	}
	//用于存储信息和回显
	var group UserGroup
	//执行sql语句并存储在group中
	Db.Select("group_id").Where("user_id = ? AND friend_id = ?", userid, friendid).First(&group)
	if group.GroupID > 0 {
		return errors.New("好友已存在于其它分组中")
	}
	return nil
}

/*
	编辑者：F,W
	功能：使用gorm创建分组
	说明：创建成功或者失败返回自定义消息状态码
*/

func CreateGroup(group *UserGroup) error {
	//主要用于检查新好友
	err := CheckGroupFriendID(group.UserID, group.FriendID)
	//friend为0是创建新分组
	if err != nil && group.FriendID != 0 {
		return err
	}
	//主要用于检查新分组
	err = CheckGroupName(group.UserID, group.GroupName)
	if err != nil {
		return err
	}
	//执行sql语句并获取错误
	err = Db.Create(&group).Error
	if err != nil {
		return err
	}
	return nil
}

/*
	编辑者：F,W
	功能：使用gorm获取指定用户的好友分组
	说明：获取成功则返回新定义结构体的相关信息
*/

func GetGroups(userid int) ([]ContactList, error) {
	//用于存储信息和回显
	var usergroups []UserGroup
	//如果传入的值不正确，则直接返回
	if userid <= 0 {
		return nil, errors.New("传入ID不正确")
	}
	//执行sql语句并返回错误信息
	err := Db.Where("user_id = ?", userid).Find(&usergroups).Error
	if err != nil {
		return nil, err
	}
	//将信息进行处理存入新切片
	groups := GetContactList(usergroups)
	return groups, nil
}

/*
	编辑者：F,W
	功能：使用gorm修改分组信息
	说明：此处是更换分组和修改备注
		通过friendId进行进一步判断
		如果friendId等于0，newGroupName不为空，表示修改整个分组，保留标记分组记录
		如果friendId大于0，newGroupName不为空，表示修改单个人的分组，如果移动分组不存在，则创建标记记录

*/

func EditGroups(data *UpdateGroup) error {
	//friend为0表示修改整个分组
	if data.UserID <= 0 || data.FriendID < 0 {
		return errors.New("用户ID不正确")
	}
	//修改整个分组
	if data.FriendID == 0 && data.NewGroupName != "" {
		//保留原分组标记记录，然后将成员转移至其它分组
		err := Db.Model(&UserGroup{}).Where("user_id = ? AND group_name = ? AND friend_id > 0", data.UserID, data.OldGroupName).Update("group_name", data.NewGroupName).Error
		if err != nil {
			return err
		}
	}
	//修改单个人所属分组
	if data.FriendID > 0 && data.NewGroupName != "" {
		err := CheckGroupName(data.UserID, data.NewGroupName)
		//如果移动的分组不存在
		if err == nil {
			err = CreateGroup(&UserGroup{GroupName: data.NewGroupName,
				UserID:   data.UserID,
				FriendID: 0})
			if err != nil {
				return err
			}
		}
		err = Db.Model(&UserGroup{}).Where("user_id = ? AND friend_id = ?", data.UserID, data.FriendID).Update("group_name", data.NewGroupName).Error
		if err != nil {
			return err
		}
	}
	//修改备注
	if data.FriendID > 0 && data.Remarks != "" {
		//直接修改检查错误即可
		err := Db.Model(&UserGroup{}).Where("user_id = ? AND friend_id = ?", data.UserID, data.FriendID).Update("remarks", data.Remarks).Error
		if err != nil {
			return err
		}
	}
	return nil
}

/*
	编辑者：F,W
	功能：使用gorm删除指定用户的好友分组，或者分组下单个好友
	说明：根据是否成功返回相关自定义状态码,如果删除整个分组，friendid应为0
*/

func DeleteGroups(group *UserGroup) error {
	log.Println(group)
	//删除整个分组
	if group.FriendID == 0 {
		//先获取列表
		var userGroups []UserGroup
		err := Db.Where("user_id = ? AND group_name = ?", group.UserID, group.GroupName).Find(&userGroups).Error
		if err != nil {
			return err
		}
		//遍历删除每个好友
		for _, item := range userGroups {
			//此处就不接收错误，直接执行
			err = Db.Where("user_id IN (?,?) AND friend_id IN (?,?)", group.UserID, item.FriendID, group.UserID, item.FriendID).Delete(&UserGroup{}).Error
			if err != nil {
				return err
			}
		}
		//删除分组标记记录,此处应该只剩一个标记可以删除了
		err = Db.Where("user_id = ？ AND group_name = ?", group.UserID, group.GroupName).Delete(&UserGroup{}).Error
		if err != nil {
			return err
		}
	} else {
		//删除单个好友
		err := Db.Where("user_id IN (?,?) AND friend_id IN (?,?)", group.UserID, group.FriendID, group.UserID, group.FriendID).Delete(&UserGroup{}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
