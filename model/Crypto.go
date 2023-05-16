package model

import (
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
	"log"
)

/*
	编辑者：Z
	功能：用户密码加密存储
	说明：在创建用户时将密码进行加密再创建进数据库
         这里给出两种方案：
         一种比较通用的bcrypt（通过视频得知比较通用）
         一种是scrypt（官方推荐的专家方案）
         但是通过实际使用，bcrypt输出长度不固定，密码长度越长，输出的越长，所有暂时选择定长的scrypt
*/

func ScryptPw(password string) string {
	//注意：密码存储长度已经设置为20
	const KeyLen = 10
	//加盐
	salt := make([]byte, 8)
	salt = []byte{12, 32, 34, 5, 7, 23, 177, 123}
	//获取加密后的内容
	HashPw, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	//因为HashPw是固定10位，base64编码固定16位，数据库存储长度为20位
	FinalPw := base64.StdEncoding.EncodeToString(HashPw)
	return FinalPw
}
