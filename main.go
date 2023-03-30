package main

import (
	"log"

	"github.com/A11Might/wechatgpt/handler"
	"github.com/A11Might/wechatgpt/helper"
	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", handler.MessageHanlder)
	r.POST("/", handler.MessageHanlder)
	// r.Run() // 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8001")
}

func Init() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	helper.DefaultConfig = &config
	helper.DefaultMessageQueue = helper.NewMessageQueue(10)
	helper.DefaultOfficialAccount = helper.NewOfficialAccount()
	helper.DefaultOpenAI = helper.NewOpenAI(5)
}
