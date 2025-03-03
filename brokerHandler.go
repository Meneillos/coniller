package main

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

const (
	EXCHANGE_NAME = "patata"
	WORK_QUEUE    = "workq"
	LOG_QUEUE     = "logq"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	workq   amqp.Queue
	logq    amqp.Queue
)

func InitBroker() {
	err := godotenv.Load()
	logOnError(err, ERROR_LOADING_ENV_FILE)
	rabbit_url, err := base64.StdEncoding.DecodeString(os.Getenv("RABBITMQ_URL"))
	logOnError(err, ERROR_DECODING_RABBITMQ_URL, true)
	conn, err = amqp.Dial(string(rabbit_url))
	logOnError(err, ERROR_CONNECTING_TO_RABBITMQ, true)
	channel, err = conn.Channel()
	logOnError(err, ERROR_OPENING_CHANNEL, true)
	initExchange()
}

func CloseBroker() {
	channel.Close()
	conn.Close()
}

func initExchange() {
	err := channel.ExchangeDeclare(
		EXCHANGE_NAME, // name
		"direct",      // type
		false,         // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	logOnError(err, ERROR_DECLARING_EXCHANGE, true)

	workq, err = channel.QueueDeclare(
		WORK_QUEUE, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	logOnError(err, ERROR_DECLARING_QUEUE, true)

	logq, err = channel.QueueDeclare(
		LOG_QUEUE, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	logOnError(err, ERROR_DECLARING_QUEUE, true)

	err = channel.QueueBind(
		workq.Name,    // queue name
		workq.Name,    // routing key
		EXCHANGE_NAME, // exchange
		false,
		nil,
	)
	logOnError(err, ERROR_BINDING_QUEUE, true)

	err = channel.QueueBind(
		logq.Name,     // queue name
		logq.Name,     // routing key
		EXCHANGE_NAME, // exchange
		false,
		nil,
	)
	logOnError(err, ERROR_BINDING_QUEUE, true)
}

func (msg Message) Publish(queue amqp.Queue) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = channel.Publish(
		EXCHANGE_NAME, // exchange
		queue.Name,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonMsg),
		})
	if err != nil {
		return err
	}
	return nil
}
