package integrations

import (
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type GetUserIntegrationEvent struct {
	Token string `json:"token"`
}

func (b *UserIntegrations) GetUserIntegrationEventHandler(d amqp091.Delivery) (any, error) {
	var data GetUserIntegrationEvent
	err := json.Unmarshal(d.Body, &data)
	if err != nil {
		panic(err)
	}

	log.Default().Println(d.ConsumerTag)

	return []byte(d.ReplyTo), nil

}
