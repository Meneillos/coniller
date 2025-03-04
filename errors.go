package main

import (
	"fmt"
	"log"
)

type Error int32

const (
	ERROR_LOADING_ENV_FILE Error = iota
	ERROR_DECODING_RABBITMQ_URL
	ERROR_CONNECTING_TO_RABBITMQ
	ERROR_OPENING_CHANNEL
	ERROR_DECLARING_EXCHANGE
	ERROR_DECLARING_QUEUE
	ERROR_BINDING_QUEUE
	ERROR_PUBLISHING_MESSAGE
	ERROR_MARSHALLING_MESSAGE
	ERROR_UNMARSHALLING_MESSAGE
	ERROR_CONSUMING_MESSAGE
	ERROR_SETTING_QOS
	ERROR_GETTING_HOSTNAME
)

func logOnError(err error, error_type Error, critical ...bool) {
	msg := fmt.Sprintf("%s | Error %d: %s", Hostname, error_type, err)
	crit := false
	if len(critical) > 0 {
		crit = critical[0]
	}
	if err != nil && !crit {
		log.Println(msg)
	}
	if err != nil && crit {
		log.Fatalln(msg)
	}
}
