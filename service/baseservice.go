package service

import "consumer/common/rabbitmq"

type (
	BaseMessageConsumerService interface {
		ConsumerMessage(message *rabbitmq.Message) error
	}
)
