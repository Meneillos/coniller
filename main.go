package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

var Hostname string

func main() {
	var err error
	Hostname, err = os.Hostname()
	logOnError(err, ERROR_GETTING_HOSTNAME, true)
	InitBroker()
	defer CloseBroker()

	var message = new(Message)
	for range 10 {
		message.Subject = fmt.Sprintf("Prueba %d", rand.Intn(1000))
		message.Body = "Pelila!!"
		err := message.Publish(workq)
		logOnError(err, ERROR_PUBLISHING_MESSAGE)
	}

	msgs, err := channel.Consume(
		workq.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	logOnError(err, ERROR_CONSUMING_MESSAGE, true)

	for m := range msgs {
		err = json.Unmarshal(m.Body, &message)
		logOnError(err, ERROR_UNMARSHALLING_MESSAGE)
		log.Printf("%s | Received a message: %s", Hostname, message.Subject)
	}
}
