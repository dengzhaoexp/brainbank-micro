package mq

import (
	"context"
	"file/config"
	"file/pkg/consts"
	fileLogger "file/pkg/utils/logger"
	"file/repositry/dao"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

var _mq *amqp.Connection

func InitRabbitMQ() {
	// 获取配置
	rConfig := config.Config.RabbitMQ

	// 获取dns
	dns := strings.Join([]string{rConfig.RabbitMQ,
		"://", rConfig.RabbitMQUser,
		":", rConfig.RabbitMQPassword,
		"@", rConfig.RabbitMQHost,
		":", rConfig.RabbitMQPort, "/"}, "")
	fileLogger.LogrusObj.Info("Load RabbitMQ configuration successfully")

	// 连接
	conn, err := amqp.Dial(dns)
	if err != nil {
		fileLogger.LogrusObj.Error("Error while getting RabbitMQ connection:", err)
		panic(err)
	}
	_mq = conn
	ch, err := _mq.Channel()
	if err != nil {
		fileLogger.LogrusObj.Error("Error while getting channel for RabbitMQ connection:", err)
		panic(err)
	}

	// 初始化基本的配置
	if err = initSetup(ch); err != nil {
		fileLogger.LogrusObj.Error("Error occurred initializing the basic RabbitMQ configuration:", err)
		panic(err)
	}
	fileLogger.LogrusObj.Info("RabbitMQ Initialization Basic Configuration Successful.")

	// 开始监听死信队列
	_ = startDeadLetterConsumer(ch)
}

func initSetup(ch *amqp.Channel) error {
	// 声明主交换机
	err := ch.ExchangeDeclare(
		consts.MainExchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 声明死信交换机
	err = ch.ExchangeDeclare(
		consts.DeadExchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 声明死信队列，如果名称相同，并不会重复声明
	_, err = ch.QueueDeclare(
		consts.DeadQueueName,
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		return err
	}

	// 绑定
	err = ch.QueueBind(consts.DeadQueueName, "dead", consts.DeadExchangeName, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func startDeadLetterConsumer(ch *amqp.Channel) error {
	// 创建死信消费
	messages, err := ch.Consume(
		consts.DeadQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fileLogger.LogrusObj.Error("Failed to get message pipeline from dead message queue:", err)
		return err
	}

	// 并发处理消息
	for i := 0; i < 5; i++ {
		// TODO:针对这里的消费失败的信息，更好的做法是设计一个消息重放机制，将处理失败的消息记录放在一个专门的错误队列当中
		go func(workerID int) {
			fileLogger.LogrusObj.Printf("Worker[%d] is started,waiting dead messages.......", workerID)
			for d := range messages { // 从通道中循环接收消息
				// 解析消息内容
				messageParts := strings.Split(string(d.Body), "|")
				if len(messageParts) != 2 {
					fileLogger.LogrusObj.Warning("Illegal message format:", string(d.Body))
					continue // 跳过无效消息
				}
				userID := messageParts[0]
				documentID := messageParts[1]
				fileLogger.LogrusObj.Printf("Dead message worker[%d] parsed userID: %s,document:%s", workerID, userID, documentID)

				// 根据ID执行删除操作
				dDao := dao.NewResourceDao(context.Background())

				// 查询文件
				document, err := dDao.GetDocumentByID(documentID)
				if err != nil {
					fileLogger.LogrusObj.Error("Dead letter queue error while fetching corresponding file from database:", err)
					continue // 跳过处理失败消息
				}

				// 更新文件状态
				document.Status = consts.FileStatusReleased
				if err = dDao.UpdateDocument(document); err != nil {
					fileLogger.LogrusObj.Error("Failed to update file status while consuming dead letter queue:", err)
					continue // 跳过处理失败消息
				}

				if err = dDao.DeleteDocumentByID(documentID); err != nil {
					fileLogger.LogrusObj.Error("Dead letter queue fails when deleting a file from the database:", err)
					continue // 跳过处理失败消息
				}

				// 消费完成，手动应答
				if err = d.Ack(false); err != nil {
					fileLogger.LogrusObj.Error("Failed to manually answer a message while consuming a dead letter queue:", err)
					continue // 跳过应答失败消息
				}

			}
		}(i)
	}
	return nil
}

func GetMQClient() *amqp.Connection {
	return _mq
}
