package main

import (
	"github.com/nsqio/go-nsq"
	"log"
)

type MessageHandler struct{}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	log.Printf("received message: %s", string(m.Body))
	return nil
}

func main() {
	// 创建消费者
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("test_topic", "test_channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// 设置消息处理器
	consumer.AddHandler(&MessageHandler{})

	// 连接到nsqd
	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}

	// 阻塞，直到退出
	<-consumer.StopChan
}
