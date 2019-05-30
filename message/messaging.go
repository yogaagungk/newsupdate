package message

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"

	"github.com/yogaagungk/newsupdate/model"

	"github.com/yogaagungk/newsupdate/news"

	"github.com/yogaagungk/newsupdate/common"
)

// Messaging merupakan contract untuk publish dan consume data dari broker
type Messaging interface {
	Listen()
	Publish(news model.News)
}

type message struct {
	amqpChannel *amqp.Channel
	service     news.Service
}

//InitDependencyMessaging digunakan untuk menginject dependency yang dibutuhkan
//pada messaging
func InitDependencyMessaging(amqpChannel *amqp.Channel, service news.Service) Messaging {
	return &message{amqpChannel, service}
}

//Listen, fungsi untuk listen dan consume message yang dipublish oleh publisher
func (message *message) Listen() {
	messageChannel, err := message.amqpChannel.Consume(
		"AddNews",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	common.HandleError(err, "Could not register consumer")

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())

		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			news := &model.News{}

			err := json.Unmarshal(d.Body, news)

			if err != nil {
				common.HandleError(err, "Error decoding Json")
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

			message.service.Save(news)
		}
	}()

}

//Publish, fungsi untuk mempublish message ke broker
func (message *message) Publish(news model.News) {
	body, err := json.Marshal(news)

	if err != nil {
		common.HandleError(err, "Error encoding Json")
	}

	err = message.amqpChannel.Publish("", "AddNews", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		common.HandleError(err, "Error publishing message: %s")
	}

}
