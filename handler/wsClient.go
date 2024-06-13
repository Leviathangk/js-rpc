package handler

/**
TODO 这个文件暂时没有使用，后期有空需要通过这个去优化整体的代码结构
*/

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type WebSocketContext struct {
	Conn    *websocket.Conn // 连接
	MsgChan chan []byte     // 消息回传 chan
	Locker  sync.Mutex      // 锁
}

func NewWebSocketContext(conn *websocket.Conn) *WebSocketContext {
	return &WebSocketContext{
		Conn:    conn,
		MsgChan: make(chan []byte),
		Locker:  sync.Mutex{},
	}
}

// Send 发送消息
func (ws *WebSocketContext) Send(msg gin.H) bool {
	ws.Locker.Lock()
	defer ws.Locker.Unlock()
	err := ws.Conn.WriteJSON(msg)
	if err != nil {
		log.Println("写入异常：", err)
		return false
	}

	return true
}

// Listen 监听消息
func (ws *WebSocketContext) Listen() {
	for {
		// 等待消息
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("关闭连接：\n")
			return
		}

		// 消息发送 chan
		ws.MsgChan <- message

		log.Printf("接收消息: %s\n", message)
	}
}
