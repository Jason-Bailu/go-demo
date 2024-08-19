package main

import (
	"github.com/nsqio/go-nsq"
	"log"
)

func main() {
	// 创建生产者
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	// 发布消息到"test_topic"
	err = producer.Publish("test_topic", []byte("hello, nsq!"))
	if err != nil {
		log.Fatal(err)
	}

	// 关闭生产者
	producer.Stop()
}
