package config

import (
	"github.com/streadway/amqp"
	"github.com/yogaagungk/newsupdate/common"
)

//Configuration struct
type Configuration struct {
	AMQPConnectionURL string
}

var config = Configuration{
	AMQPConnectionURL: "amqp://guest:12345@localhost:5672/",
}

// InitialChannel , konfigurasi dan open connection ke broker,
// Disini broker yang digunakan adalah AMPQ RabbitMQ
func InitialChannel() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	connection, err := amqp.Dial(config.AMQPConnectionURL)
	common.HandleError(err, "Can't connect to AMPQ")

	amqpChannel, err := connection.Channel()
	common.HandleError(err, "Can't create AMPQ Channel")

	queue, err := amqpChannel.QueueDeclare("AddNews", true, false, false, false, nil)
	common.HandleError(err, "Could not declare queue 'AddNews'")

	err = amqpChannel.Qos(1, 0, false)
	common.HandleError(err, "Could not configure QoS")

	return connection, amqpChannel, &queue
}
