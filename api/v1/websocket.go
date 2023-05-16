package v1

import (
	"chat-online/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
	编辑者：F
	功能：websocket升级器
	说明：设置一些内容
*/

var upGrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		getID := r.URL.Query().Get("id")
		if getID == "" {
			return false
		}
		return true
	},
}

/*
	编辑者：F
	功能：websocket
	说明：此处未本系统的核心，需要将所有登陆后的操作都集中在这里
         在登陆成功后主动建立ws连接，并进行相关的交互，传递userid、token和operation
*/

func WebSocket(c *gin.Context) {
	//用于存储聊太信息
	var ws model.WsCarrier
	getID := c.Query("id")
	getIDint, err := strconv.Atoi(getID)
	if err != nil {
		return
	}
	//通过升级后的升级器得到链接
	conn, err2 := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err2 != nil {
		log.Println("获取连接失败:", err2)
		return
	} else {
		log.Println("connect")
	}
	defer conn.Close()
	//得到连接后，就可以开始读写数据了
	model.UserList[getIDint] = model.OnlineUser{
		ID:     getIDint,
		Socket: conn,
	}
	for {
		err := conn.ReadJSON(&ws)
		if err != nil {
			return
		}
		model.HandledWebSocket(&ws)
	}
}
