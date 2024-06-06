package mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RecoveryConsume(ch *amqp.Channel, documentID string) (bool, error) {
	// 确定队列名
	queueName := fmt.Sprintf("document_queue_%s", documentID)

	// 消息消费
	message, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return false, err
	}

	// 获取消息
	_, ok := <-message

	// 消息消费后销毁该队列名
	if _, err = ch.QueueDelete(queueName, false, false, false); err != nil {
		return false, err
	}

	return ok, nil
}
