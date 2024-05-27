package integrations

import (
	"services/internal/storage/services"
	"services/pkg/eventbus/rabbitmq"
)

type StorageIntegrations struct {
	broker       *rabbitmq.RabbitMQ
	service_user *services.StorageService
}

func NewIntegrations(broker *rabbitmq.RabbitMQ, user *services.StorageService) (*StorageIntegrations, error) {

	return &StorageIntegrations{
		broker:       broker,
		service_user: user,
	}, nil
}
