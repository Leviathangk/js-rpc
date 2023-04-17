package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete 删除 rpc
func Delete(c *gin.Context) {
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

	msg["type"] = TypeDelete

	// 等待通道消息
	WaitChan(c, msg)
}

// Delete 删除 rpc
func DeleteByUUID(c *gin.Context) {
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

	msg["type"] = TypeDelete

	// 等待通道消息
	WaitChanByUUID(c, msg)
}

// Delete 删除 rpc
func DeleteMore(c *gin.Context) {
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

	msg["type"] = TypeDelete

	// 等待通道消息
	WaitChanMore(c, msg)
}
