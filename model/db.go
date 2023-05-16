package model

import (
	"chat-online/utils"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	err error
)

func InitDb() {
	//目标数据库内容
	dstDb := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.DbUser,
		utils.DbPassword,
		utils.DbHost,
		utils.DbPort,
		utils.DbName)
	//设置内部禁用自动创建外键约束
	Db, err = gorm.Open(mysql.Open(dstDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Printf("连接数据库失败：%v\n", err)
	}
	//自动迁移
	Db.AutoMigrate(&User{}, &UserGroup{}, &ChatLog{}, &AddRequest{}, &EmailVerification{})
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err2 := Db.DB()
	if err2 != nil {
		fmt.Printf("获取通用数据库对象失败：%v\n", err2)
	}
	//defer sqlDB.Close()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}
