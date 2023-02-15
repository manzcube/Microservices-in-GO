package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Check RabbitMQ documentation
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // topic,
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguments?

	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable(get rid after done)
		false, //delete when unsed?
		true,  // exclusive
		false, // no-wait?
		nil,   // arguments?
	)
}
