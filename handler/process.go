package handler

import (
	"encoding/json"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/gin-gonic/gin"
)

// SendJson ws 情况下发送 json 消息
func (c *Client) SendJson(m gin.H) {
	err := c.Conn.WriteJSON(m)
	if err != nil {
		glog.Warnf("【%s】发送消息失败：%s\n", c.UUID, err.Error())
	}
}

// Stop 关闭并移除客户端
func (c *Client) Stop() {
	manager.RemoveClient(c.UUID)
	c.Conn.Close()
	glog.Debugf("【%s】清理资源...\n", c.UUID)
}

// ListenMsg 监听消息
func (c *Client) ListenMsg() {
	defer c.Stop()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			glog.Errorf("【%s】关闭连接：%v\n", c.UUID, err)
			break
		}

		glog.Debugf("【%s】收到消息：%s\n", c.UUID, msg)
		go c.ProcessMsg(msg)

		// 发送消息
		// c.Conn.WriteMessage(mt, []byte("你好"))
	}
}

// ProcessMsg 处理消息
func (c *Client) ProcessMsg(m []byte) {
	msg := new(Message)
	err := json.Unmarshal(m, msg)
	if err != nil {
		glog.Warnf("json 解析失败：%s\n", string(m))
		return
	}

	// 针对消息类型做不同处理
	switch msg.Type {
	case TypeOpen:
		domain := msg.Msg["domain"].(string)
		c.Domain = domain
		manager.Domains[domain] = append(manager.Domains[domain], c)
		glog.Debugf("向域名 %s 下存入机器 %s，当前数量：%d\n", domain, c.UUID, len(manager.Domains[domain]))
	default:
		if msgContext, ok := managerMsg[msg.EventId]; ok {
			func() {
				msgContext.Locker.Lock()
				defer msgContext.Locker.Unlock()
				if !msgContext.IsStop {
					msgContext.IsStop = true
					msgContext.MsgChan <- msg.Msg
				}
				delete(managerMsg, msg.EventId)
			}()
		} else {
			glog.Warnf("消息通道不存在：%s\n", msg.EventId)
		}
	}
}

// func (c *Client)SendJsonMsg()
