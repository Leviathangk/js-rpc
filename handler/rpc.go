package handler

import (
	"net/http"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// 最大超时执行时间
	MaxRunTime = 5 * 60
)

var (
	// 统一管理 web 端
	manager = NewManager()

	// 统一管理消息端
	managerMsg = map[string]*MsgContext{} // 事件 id:chan

	// 将 http 升级到 ws
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 设置允许跨域
		},
	}
)

func Rpc(c *gin.Context) {
	var err error
	var wsConn *websocket.Conn

	// 获取客户端的 uuid
	clientUUID := c.Query("uuid")
	if len(clientUUID) == 0 {
		clientUUID = CreateUUID()
	}

	// 升级为 ws 连接
	wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": true,
			"msg":     "连接失败：" + err.Error(),
		})
		return
	}

	// 保存客户端
	client := manager.AddClient(wsConn, clientUUID)
	glog.Debugf("连接建立：%s -> %s\n", wsConn.RemoteAddr().String(), client.UUID)

	// 发送初始消息
	client.SendJson(gin.H{
		"uuid": client.UUID,
		"type": TypeOpen,
	})

	// 监听消息
	client.ListenMsg()
}
