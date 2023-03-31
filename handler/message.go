package handler

import (
	"fmt"

	"github.com/A11Might/wechatgpt/biz"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func MessageHanlder(c *gin.Context) {
	server := biz.DefaultMessageService.GetOfficialAccountService(c)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		reply, err := biz.DefaultMessageService.ProcessMessage(c, msg)
		if err != nil {
			reply = message.NewText("处理错误请重试")
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: reply}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	fmt.Printf("got:%+v\n", server.ResponseMsg)
	_ = server.Send()
}
