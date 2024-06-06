package mq

import (
	"context"
	"file/pkg/consts"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitUserQueueAndBind(ch *amqp.Channel, documentID string) error {
	// 定义用户队列名称
	queueName := fmt.Sprintf("document_queue_%s", documentID)

	// 声明用户队列
	if _, err := ch.QueueDeclare(
		queueName,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		amqp.Table{
			"x-dead-letter-routing-key": "dead",                  // routing-key
			"x-dead-letter-exchange":    consts.DeadExchangeName, // 死信交换机
			"x-message-ttl":             consts.MessageTTL},      // TTL
	); err != nil {
		return err
	}

	// 绑定队列
	rk := fmt.Sprintf("%s", documentID)
	if err := ch.QueueBind(queueName, rk, consts.MainExchangeName, false, nil); err != nil {
		return err
	}

	return nil
}

func PublishMessage(ch *amqp.Channel, userId string, documentID string) error {
	// 构建信息体
	message := fmt.Sprintf("%s|%s", userId, documentID)

	// 发布信息
	ctx := context.Background()
	err := ch.PublishWithContext(ctx,
		consts.MainExchangeName,
		documentID,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
