package helper

import (
	"errors"
	"log"
	"strconv"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type MessageQueue struct {
	q      Queue[*message.MixMessage]
	loader *SingleFlightGroup
}

func NewMessageQueue(capacity int) *MessageQueue {
	return &MessageQueue{
		q:      NewQueue[*message.MixMessage](capacity),
		loader: &SingleFlightGroup{},
	}
}

func (ms *MessageQueue) AddMessage(msg *message.MixMessage) error {
	if !ms.q.TryPush(msg) {
		return errors.New("排队处理的消息太多了，请稍后再试")
	}
	return nil
}

func (ms *MessageQueue) ProcessMessage(request *message.MixMessage) (*message.Text, error) {
	log.Println("processing message")
	replyIface, err := ms.loader.Do(strconv.Itoa(int(request.MsgID)), func() (interface{}, error) {
		text, err := DefaultOpenAI.Chat(request.Content)
		if err != nil {
			log.Printf("process message error, err:%+v", err)
			text = "消息处理出错了"
		}
		return message.NewText(text), nil
	})
	return replyIface.(*message.Text), err
}

var DefaultMessageQueue *MessageQueue
