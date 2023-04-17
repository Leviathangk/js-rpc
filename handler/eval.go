package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Eval 运行 rpc
func Eval(c *gin.Context) {
	fields := []string{"domain", "js"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺失字段：" + missField,
		})
		return
	}

	msg["type"] = TypeEval

	// 等待通道消息
	WaitChan(c, msg)
}

// Eval 运行 rpc
func EvalByUUID(c *gin.Context) {
	fields := []string{"uuid", "js"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺失字段：" + missField,
		})
		return
	}

	msg["type"] = TypeEval

	// 等待通道消息
	WaitChanByUUID(c, msg)
}
