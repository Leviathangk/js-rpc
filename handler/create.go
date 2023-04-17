package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create 通过 domain 给随机该域名客户端创建 rpc
func Create(c *gin.Context) {
	fields := []string{"domain", "funcName", "funcBody"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺失字段：" + missField,
		})
		return
	}

	msg["type"] = TypeCreate

	// 等待通道消息
	WaitChan(c, msg)
}

// CreateByUUID 通过 uuid 创建 rpc
func CreateByUUID(c *gin.Context) {
	fields := []string{"uuid", "funcName", "funcBody"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺失字段：" + missField,
		})
		return
	}

	msg["type"] = TypeCreate

	// 等待通道消息
	WaitChanByUUID(c, msg)
}

// CreateMore 通过 domain 给所有该域名客户端创建 rpc
func CreateMore(c *gin.Context) {
	fields := []string{"domain", "funcName", "funcBody"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "缺失字段：" + missField})
		return
	}

	msg["type"] = TypeCreate

	// 等待通道消息
	WaitChanMore(c, msg)
}
