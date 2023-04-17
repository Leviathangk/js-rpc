package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowDomainClients 查看指定域名下的所有机器的 id
func ShowDomainClients(c *gin.Context) {
	fields := []string{"domain"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "缺失字段：" + missField})
		return
	}

	// 查看本地机器
	domain := msg["domain"].(string)
	clients, exists := manager.GetClientsByDomain(domain)
	if !exists {
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": gin.H{"total": 0, "clients": []string{}}})
		return
	}

	uuids := []string{}
	for _, c := range clients {
		uuids = append(uuids, c.UUID)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": gin.H{"total": len(uuids), "clients": uuids}})
}

// ShowClientFunctions 查看指定机器下的函数列表
func ShowClientFunctions(c *gin.Context) {
	fields := []string{"uuid"}

	// 获取、检查、构造消息
	missField, msg := fieldsCheck(c, fields)
	if missField != "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "缺失字段：" + missField})
		return
	}

	msg["type"] = TypeShow

	WaitChanByUUID(c, msg)
}
