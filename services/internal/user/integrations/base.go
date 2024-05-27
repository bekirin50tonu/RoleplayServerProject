package integrations

import (
	"services/internal/user/services"
	"services/pkg/eventbus/rabbitmq"
)

type UserIntegrations struct {
	broker       *rabbitmq.RabbitMQ
	service_user *services.UserService
}

func NewIntegrations(broker *rabbitmq.RabbitMQ, user *services.UserService) (*UserIntegrations, error) {

	return &UserIntegrations{
		broker:       broker,
		service_user: user,
	}, nil
}
