package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"services/pkg/common/response"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL            string
	ExchangeName   string
	ExchangeSuffix string
	ExchangePrefix string
	AcK            bool
	PrefetchCount  int
}

type IntegrationEvent struct {
	Data interface{} `json:"data"`
}

type RabbitMQ struct {
	conn   *amqp.Connection
	memory map[string]*amqp.Channel
	config RabbitMQConfig
	mutex  sync.Locker
}

type IRabbitMQ interface {
	Publish()
	AddSubscription()
	RemoveSubscription()
}

func NewRabbitMQ(config RabbitMQConfig) (*RabbitMQ, error) {

	// Context oluşturun ve bir süre sınırlayın (10 saniye)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Context kullanım sona erdiğinde iptal et

	// RabbitMQ bağlantısı kurma işlemi
	var conn *amqp.Connection
	var err error

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Bağlantı süresi aşıldı.")
			return nil, errors.New("Connection Timeout.")
		default:
			conn, err = amqp.Dial(config.URL)
			if err == nil {
				fmt.Println(fmt.Sprintf("RabbitMQ bağlantısı başarılı.State:%v", !conn.IsClosed()))
				return &RabbitMQ{
					conn:   conn,
					memory: make(map[string]*amqp.Channel),
					config: config,
				}, nil
			}
			fmt.Printf("Bağlantı kurma hatası: %v. Yeniden denenecek...\n", err)
			time.Sleep(1 * time.Second) // Bağlantı hatası durumunda 1 saniye bekleme
		}
	}

}

func (a *RabbitMQ) Close() {

	err := a.conn.Close()
	if err != nil {
		log.Fatal("%w", err)
	}
	for _, channel := range a.memory {
		err = channel.Close()
		if err != nil {
			log.Fatal("%w", err)
		}
	}

}

func (a *RabbitMQ) AddSubscription(test interface{}, handler interface{}) (any, error) {
	exchange := a.config.ExchangeName
	queue := reflect.TypeOf(test).Name()
	queue = a.config.ExchangeSuffix + "." + strings.Replace(queue, a.config.ExchangePrefix, "", -1)

	_, ok := a.memory[queue]
	if ok {
		return nil, errors.New("Queue is Already Declared.")
	}

	channel, err := a.createChannel(a.config.PrefetchCount)

	if err != nil {
		fmt.Printf("State Channel:%v", !channel.IsClosed())
		return nil, err
	}

	err = a.declareExchange(channel, exchange, "direct")
	if err != nil {
		return nil, err
	}

	_, err = a.declareQueue(channel, queue)
	if err != nil {
		return nil, err
	}

	_, err = a.declareDeadLetterQueue(channel, queue)
	if err != nil {
		return nil, err
	}

	err = a.bindQueue(channel, queue, exchange)
	if err != nil {
		return nil, err
	}

	_, err = a.initializeCostumer(channel, queue, test, handler)
	if err != nil {
		return nil, err
	}
	a.memory[queue] = channel
	fmt.Println("Subscribed:", queue)
	return nil, nil
}

func (a *RabbitMQ) RemoveSubscription(queue string) {

}

func (a *RabbitMQ) Publish(exchange string, routingKey string, data interface{}) error {
	channel, err := a.createChannel(a.config.PrefetchCount)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	encodedData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Timestamp:     time.Now(),
		CorrelationId: uuid.NewString(),
		Body:          encodedData,
		ReplyTo:       "amq.rabbitmq.reply-to",
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *RabbitMQ) PublishAndWaitForReply(exchange string, routingKey string, data interface{}) (chan amqp.Delivery, error) {

	channel, err := a.createChannel(a.config.PrefetchCount)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	encodedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume("amq.rabbitmq.reply-to", routingKey, true, false, false, false, nil)
	resp := make(chan amqp.Delivery)
	go func(deliveries <-chan amqp.Delivery) {
		for d := range deliveries {
			resp <- d
		}
	}(msgs)

	if err != nil {
		return nil, err
	}
	fmt.Println("geas", routingKey)
	err = channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Timestamp:     time.Now(),
		CorrelationId: uuid.NewString(),
		Body:          encodedData,
		ReplyTo:       "amq.rabbitmq.reply-to",
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *RabbitMQ) createChannel(prefetchCount int) (*amqp.Channel, error) {
	channel, err := a.conn.Channel()
	if err != nil {
		return nil, err
	}
	err = channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (a *RabbitMQ) declareDeadLetterQueue(channel *amqp.Channel, queueName string) (*amqp.Queue, error) {

	queue, err := channel.QueueDeclare(queueName+".deadLetter", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &queue, err
}

func (a *RabbitMQ) declareQueue(channel *amqp.Channel, queueName string) (*amqp.Queue, error) {
	arg := amqp.Table{"x-dead-letter-exchange": a.config.ExchangeName,
		"x-dead-letter-routing-key": queueName + ".deadLetter"}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, arg)
	if err != nil {
		return nil, err
	}

	return &queue, nil

}

func (a *RabbitMQ) declareExchange(channel *amqp.Channel, name string, kind string) error {
	err := channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
	return err
}

func (a *RabbitMQ) bindQueue(channel *amqp.Channel, queue string, exchange string) error {
	err := channel.QueueBind(queue, queue, a.config.ExchangeName, false, nil)
	return err
}

func (a *RabbitMQ) initializeCostumer(channel *amqp.Channel, queue string, param interface{}, handler interface{}) (any, error) {
	deliveries, err := channel.Consume(queue, queue, a.config.AcK, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	// Start goroutine and consume coming data.
	go func(delivery <-chan amqp.Delivery) {
		for d := range delivery {

			// Call Function.
			r, err := handler.(func(*amqp.Delivery) (any, error))(&d)

			// If corrupted any error in calling function.
			if err != nil {
				errorMessages := make([]response.ErrorMessages, 0)

				// Has Validation Error?
				validationErrs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, v := range validationErrs {

						message := response.ErrorMessages{
							Code:    rand.Int(),
							Message: fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", v.Field(), v.Value(), v.Tag()),
						}
						errorMessages = append(errorMessages, message)
					}

				} else {
					// Has Default Error?
					errorMessages = append(errorMessages, response.ErrorMessages{
						Code:    rand.Int(),
						Message: err.Error(),
					})
				}
				response := response.ResponseWithError{
					Success: false,
					Errors:  errorMessages,
				}
				resp, err := json.Marshal(response)
				if err != nil {
					log.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				defer cancel()
				err = channel.PublishWithContext(ctx, "", d.ReplyTo, false, false, amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: uuid.NewString(),
					Body:          resp,
				})
			}

			// If response has any data. Means of have a data and must send function response to receiver.
			if r != nil {
				// Convert Json data from coming data.
				respJson, err := json.Marshal(r)
				// Convert Error.
				if err != nil {
					log.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				defer cancel()
				err = channel.PublishWithContext(ctx, "", d.ReplyTo, false, false, amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: uuid.NewString(),
					Body:          respJson,
				})

				if err != nil {
					log.Fatal(err)
				}
			}
			// Default config ack value is false.
			if !a.config.AcK {
				err = d.Ack(false)
				if err != nil {
					log.Fatal(err)
				}
			}

			// end

		}

	}(deliveries)
	return true, nil
}
