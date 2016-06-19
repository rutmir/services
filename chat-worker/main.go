package main

import (
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/rutmir/services/core/log"
	dto "github.com/rutmir/services/entities/dto/v2"
	//"github.com/rutmir/services/core/memcache"
)

const (
	atPrefix = "at_"
	rtPrefix = "rt_"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main() {
	log.Info("Initialize chat worker")

	host := os.Getenv("AMQP_HOST")
	if len(host) == 0 {
		log.Fatal("AMQP error: Required to set 'AMQP_HOST' environment")
		return
	}

	port := os.Getenv("AMQP_PORT")
	if len(port) == 0 {
		log.Fatal("AMQP error: Required to set 'AMQP_PORT' environment")
		return
	}

	user := os.Getenv("AMQP_USERNAME")
	if len(user) == 0 {
		log.Fatal("AMQP error: Required to set 'AMQP_USERNAME' environment")
		return
	}

	password := os.Getenv("AMQP_PASS")
	if len(password) == 0 {
		log.Fatal("AMQP error: Required to set 'AMQP_PASS' environment")
		return
	}

	amqpServerURL := "amqp://" + user + ":" + password + "@" + host + ":" + port + "/"

	log.Info(amqpServerURL)

	conn, err := amqp.Dial(amqpServerURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chat_worker", // name
		false,         // durable
		false,         // autoDelete
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			im := new(dto.InternalMessage)
			err := proto.Unmarshal(d.Body, im)
			failOnError(err, "Failed to decode body to InternalMessage")

			log.Info("Work on message: %s, action: %s", d.CorrelationId, im.Header.Action)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(im.Header.Action),
				})
			failOnError(err, "Failed to publish a message")
			d.Ack(false)
		}
	}()

	log.Info(" [*] Awaiting RPC requests")
	<-forever
}
