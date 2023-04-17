package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Run 运行 rpc
func Run(c *gin.Context) {
	fields := []string{"domain", "funcName"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺失字段：" + missField,
		})
		return
	}

	msg["type"] = TypeRun

	// 等待通道消息
	WaitChan(c, msg)
}

// Run 运行 rpc
func RunByUUID(c *gin.Context) {
	fields := []string{"uuid", "funcName"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺失字段：" + missField,
		})
		return
	}

	msg["type"] = TypeRun

	// 等待通道消息
	WaitChanByUUID(c, msg)
}
