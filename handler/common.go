package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/gin-gonic/gin"
)

// fieldsCheck 字段检查，有问题就会返回字段
func fieldsCheck(c *gin.Context, fields []string) (string, gin.H) {
	fieldMap := gin.H{}

	for _, field := range fields {
		if fieldVal := c.PostForm(field); fieldVal != "" {
			fieldMap[field] = fieldVal
		} else {
			return field, nil
		}
	}

	return "", fieldMap
}

// getTimeOut 获取超时时间
func getTimeOut(c *gin.Context) (timeout int) {
	var err error

	// 获取超时时间
	timeOutStr := c.PostForm("timeout")
	if timeOutStr != "" {
		timeout, err = strconv.Atoi(timeOutStr)
		if err != nil {
			timeout = MaxRunTime
		}
	} else {
		timeout = MaxRunTime
	}

	return
}

// WaitChan 等待通道消息并发送
func WaitChan(c *gin.Context, msg gin.H) {
	domain := msg["domain"].(string)

	// 获取超时时间
	timeout := getTimeOut(c)

	// 创建函数
	client, exists := manager.GetClientByDomain(domain)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "域名 " + domain + " 下无机器存在"})
		return
	}

	// 处理消息
	send(msg, []*Client{client}, timeout, func(backMsg gin.H) {
		c.JSON(http.StatusOK, backMsg["msg"].(gin.H))
	})
}

// WaitChanByUUID 等待通道消息并发送
func WaitChanByUUID(c *gin.Context, msg gin.H) {
	uuid := msg["uuid"].(string)

	// 获取超时时间
	timeout := getTimeOut(c)

	// 创建函数
	client, exists := manager.GetClientByUUID(uuid)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "客户端不存在：" + uuid})
		return
	}

	// 处理消息
	send(msg, []*Client{client}, timeout, func(backMsg gin.H) {
		c.JSON(http.StatusOK, backMsg["msg"].(gin.H))
	})
}

// WaitChan 等待通道消息并发送
func WaitChanMore(c *gin.Context, msg gin.H) {
	domain := msg["domain"].(string)

	// 获取超时时间
	timeout := getTimeOut(c)

	// 创建函数
	clients, exists := manager.GetClientsByDomain(domain)
	if !exists {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "域名 " + domain + " 下无机器存在"})
		return
	}

	// 处理消息
	success := []gin.H{}
	failed := []gin.H{}

	send(msg, clients, timeout, func(backMsg gin.H) {
		if backMsg["success"].(bool) {
			success = append(success, backMsg["msg"].(gin.H))
		} else {
			failed = append(failed, backMsg["msg"].(gin.H))
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg": gin.H{
			"success":      success,
			"failed":       failed,
			"total":        len(clients),
			"successTotal": len(success),
			"failedTotal":  len(failed),
		},
	})
}

// send 发送消息
func send(msg gin.H, clients []*Client, timeout int, callback func(backMsg gin.H)) {
	// 遍历注入
	backChan := make(chan gin.H)
	defer close(backChan)
	for _, client := range clients {
		go wait(backChan, client, msg, timeout)
	}

	// 获取结果
	for i := 0; i < len(clients); i++ {
		v := <-backChan
		callback(v)
	}
}

// wait 等待消息处理
func wait(backChan chan gin.H, client *Client, msg map[string]any, timeout int) {
	// 创建消息回发通道、事件 id
	msgChan := make(chan gin.H, 1) // 消息回传通道
	defer close(msgChan)

	// 事件上下文
	eventId := CreateUUID()
	msg["eventId"] = eventId // 事件 id
	msgContext := &MsgContext{
		MsgChan: msgChan,
	}
	managerMsg.Set(eventId, msgContext) // 存储事件 eventId:通道

	defer func() {
		if r := recover(); r != nil {
			glog.Errorln("Recovered in wait", r)

			backChan <- gin.H{
				"success": false,
				"msg": gin.H{
					"success": false,
					"msg":     "运行失败，发生 panic 错误！",
					"uuid":    client.UUID,
				},
			}
		}
	}()

	// 发出消息
	client.Locker.Lock()
	err := client.Conn.WriteJSON(msg)
	client.Locker.Unlock()
	if err != nil {
		backChan <- gin.H{
			"success": false,
			"msg": gin.H{
				"success": false,
				"msg":     "运行失败：" + err.Error(),
				"uuid":    client.UUID,
			},
		}
		return
	}

	// 创建超时信道
	timeChan := make(chan bool, 1)
	msgContext.TimeoutChan = timeChan
	defer close(timeChan)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			timeout -= 1
			if msgContext.IsStop {
				break
			} else if timeout == 0 {
				func() {
					msgContext.Locker.Lock()
					defer msgContext.Locker.Unlock()
					if !msgContext.IsStop {
						msgContext.IsStop = true
						msgContext.TimeoutChan <- true
					}
					managerMsg.Delete(eventId)
				}()
				break
			}
		}
	}()

	// 等待消息并回发
	select {
	case v := <-msgChan:
		v["uuid"] = client.UUID
		backChan <- gin.H{
			"success": true,
			"msg":     v,
		}
	case <-timeChan:
		backChan <- gin.H{
			"success": false,
			"msg": gin.H{
				"success": false,
				"msg":     "运行超时",
				"uuid":    client.UUID,
			},
		}
	}
}
