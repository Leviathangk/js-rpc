package handler

import (
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	TypeOpen = iota
	TypeCreate
	TypeShow
	TypeDelete
	TypeRun
	TypeEval
)

type MsgContext struct {
	MsgChan     chan gin.H // 消息回传 chan
	TimeoutChan chan bool  // 超时 chan
	IsStop      bool       // 是否已停止
	Locker      sync.Mutex // 锁
}

type Message struct {
	Type    int    `json:"type"`    // 消息类型
	Msg     gin.H  `json:"msg"`     // 消息
	EventId string `json:"eventId"` // 事件 id
}
