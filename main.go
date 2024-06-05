package main

import (
	"bufio"
	"fmt"
	"js-rpc/handler"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Leviathangk/go-glog/glog"
	"github.com/gin-gonic/gin"
)

var (
	serverAddr = "127.0.0.1:8080"
)

func index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "现在是：" + time.Now().String(),
	})
}

func getAddr() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("默认启动 addr：%s，输入 out 可随时退出使用默认\n", serverAddr)
		fmt.Print("请输入 ip:port：")
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		text = strings.TrimSpace(text)
		if text == "out" {
			break
		} else if text != "" {
			serverAddr = text
			break
		}
	}
}

func main() {
	// getAddr()
	router := gin.Default()
	router.GET("/", index)          // 首页：http://127.0.0.1:8080/
	router.GET("/rpc", handler.Rpc) // rpc：127.0.0.1:8080/rpc

	// 通用方法
	router.POST("/rpc/show_domain_clients", handler.ShowDomainClients)     // 查看指定域名下的所有机器的 uuid
	router.POST("/rpc/show_client_functions", handler.ShowClientFunctions) // 查看指定机器下的所有函数

	// 针对单个随机机器：通过域名
	router.POST("/rpc/domain/create", handler.Create) // 创建函数
	router.POST("/rpc/domain/delete", handler.Delete) // 删除函数
	router.POST("/rpc/domain/run", handler.Run)       // 执行函数
	router.POST("/rpc/domain/eval", handler.Eval)     // 执行 eval

	// 针对单个指定机器：通过 uuid
	router.POST("/rpc/uuid/create", handler.CreateByUUID) // 创建函数
	router.POST("/rpc/uuid/delete", handler.DeleteByUUID) // 删除函数
	router.POST("/rpc/uuid/run", handler.RunByUUID)       // 执行函数
	router.POST("/rpc/uuid/eval", handler.EvalByUUID)     // 执行 eval

	// 针对多个机器：多个机器同时执行：通过域名
	router.POST("/rpc/domain/more/create", handler.CreateMore) // 创建函数
	router.POST("/rpc/domain/more/delete", handler.DeleteMore) // 删除函数

	glog.Debugf("服务已启动：%s\n", serverAddr)
	router.Run(serverAddr)
}
