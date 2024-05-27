package main

import (
	"encoding/json"
	"fmt"
	"services/pkg/eventbus/rabbitmq"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestBroker(t *testing.T) {
	fmt.Println("TestBroker Starting...")
	broker, err := rabbitmq.NewRabbitMQ(rabbitmq.RabbitMQConfig{PrefetchCount: 10, URL: "amqp://localhost:5672", ExchangeName: "Replica.User", AcK: false, ExchangePrefix: "IntegrationEvent", ExchangeSuffix: "User"})
	if err != nil {
		t.Fatal(err)
	}

	type AddUserIntegrationEvent struct {
		Name string `json:"name"`
	}

	AddUserIntegrationEventHandler := func(d *amqp091.Delivery) (any, error) {
		var data AddUserIntegrationEvent
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	broker.AddSubscription(AddUserIntegrationEvent{}, AddUserIntegrationEventHandler)
	msgs, err := broker.PublishAndWaitForReply("Replica.User", "User.AddUser", AddUserIntegrationEvent{Name: "Bu bir testtir"})
	resp := <-msgs
	var data AddUserIntegrationEvent
	err = json.Unmarshal(resp.Body, &data)

	if err != nil {
		t.Error("Parse Error")
	}
	fmt.Println("gelen", data.Name)
	if err != nil {
		t.Error("Test Başarısız.")
	}
	t.Log("Pass")
	/* else {
		var resp AddUserIntegrationEvent

		err = json.Unmarshal(response.Body, &resp)
		if err != nil {
			fmt.Errorf("%v", err)
		}

		if resp.Name != "Bu bir testtir" {
			fmt.Errorf("Test Başarısız: %v", resp)
		}

		fmt.Println("Test Başarılı", resp)
	} */

	defer broker.Close()

}
