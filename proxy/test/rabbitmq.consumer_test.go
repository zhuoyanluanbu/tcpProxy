package test

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"testing"
)

func TestRabbitmqConsumer(t *testing.T) {
	var (
		uri = "amqp://hyc:hyc123@192.168.0.13:5678"
	)

	connection, err := amqp.Dial(uri)
	if err != nil {
		fmt.Errorf("Dial: %s", err)
	}

	channel, _ := connection.Channel()

	delivery, err := channel.Consume(
		"delay.queue",    // name
		"delay-consumer", // consumerTag,
		false,            // noAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // arguments
	)
	if err != nil {
		fmt.Errorf("Queue Consume: %s", err)
	}

	for d := range delivery {
		log.Printf(
			"deliveries2 got from [exchange->%s],[key->%s],[content->%q]",
			d.Exchange,
			d.RoutingKey,
			d.Body,
		)
		d.Ack(false)
	}

}
