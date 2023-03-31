package biz

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/A11Might/wechatgpt/helper"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"golang.org/x/sync/singleflight"
)

type MessageService struct {
	oa  *helper.OfficialAccount
	oai *helper.OpenAI

	c      *cache.Cache
	loader *singleflight.Group
}

func NewMessageService(contextSize int) *MessageService {
	return &MessageService{
		oa:     helper.NewOfficialAccount(),
		oai:    helper.NewOpenAI(contextSize),
		c:      cache.New(15*time.Second, 10*time.Minute),
		loader: &singleflight.Group{},
	}
}

func (ms *MessageService) ProcessMessage(ctx context.Context, request *message.MixMessage) (*message.Text, error) {
	log.Println("processing message")
	cacheKey := strconv.Itoa(int(request.MsgID))
	if replyIface, ok := ms.c.Get(cacheKey); ok {
		return replyIface.(*message.Text), nil
	}

	replyIface, err, _ := ms.loader.Do(cacheKey, func() (interface{}, error) {
		text, err := ms.oai.Chat(ctx, request.Content)
		if err != nil {
			log.Printf("process message error, err:%+v", err)
			text = "消息处理出错了"
		}
		reply := message.NewText(text)
		ms.c.Set(cacheKey, reply, cache.DefaultExpiration)
		return reply, nil
	})
	return replyIface.(*message.Text), err
}

func (ms *MessageService) GetOfficialAccountService(c *gin.Context) *server.Server {
	return ms.oa.GetServer(c)
}

var DefaultMessageService *MessageService
