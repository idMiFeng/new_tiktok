package mq

import (
	"fmt"
	"video_service/config"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var (
	Producer rocketmq.Producer
)

func Init() (err error) {
	Producer, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.Conf.RocketMQConfig.Addr})),
		//producer.WithNsResolver(primitive.NewPassthroughResolver(endPoint)),
		producer.WithRetry(2),
		producer.WithGroupName(config.Conf.RocketMQConfig.GroupID),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = Producer.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	return nil
}

func Exit() error {
	err := Producer.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
	return err
}
