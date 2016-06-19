package main

import (
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/rutmir/services/core/log"
	"github.com/rutmir/services/core/memcache"
	dto "github.com/rutmir/services/entities/dto/v2"
)

const (
	atPrefix = "at_"
	rtPrefix = "rt_"
)

var memCtrl memcache.MemCache
var channel *amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func responseInternalError(d *amqp.Delivery) {
	log.Info("Bad request response")
	h := new(dto.Header)
	h.Action = dto.Action_Result
	h.Timestamp = time.Now().UnixNano()

	result := new(dto.Result)
	result.Code = 500
	result.Result = "Error"
	result.Message = "Internal server error"

	if data, err := proto.Marshal(result); err != nil {
		log.Err(err)
	} else {
		msg := new(dto.InternalMessage)
		msg.Header = h
		msg.Body = data

		if data, err := proto.Marshal(msg); err != nil {
			log.Err(err)
		} else {
			err = channel.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/x-protobuf",
					CorrelationId: d.CorrelationId,
					Body:          data,
				})
			failOnError(err, "Failed to publish a message")
		}
	}
}

func main() {
	log.Info("Initialize chat worker")

	var err error
	memCtrl, err = memcache.GetLocalInstance("memcached", "test")
	if err != nil {
		log.Fatal(err)
		return
	}

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

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"chat_worker", // name
		false,         // durable
		false,         // autoDelete
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := channel.Consume(
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

			mem, err := memCtrl.Get(atPrefix + d.ReplyTo)
			if err != nil {
				log.Err(err)
				responseInternalError(&d)
				continue
			}

			auth := new(dto.AuthTokenMem)
			if err := proto.Unmarshal(mem.Value, auth); err != nil {
				log.Err(err)
				responseInternalError(&d)
				continue
			}

			log.Info("Work on message: %s, action: %s, for: %s", d.CorrelationId, im.Header.Action, auth.ProfileID.Hex())

			err = channel.Publish(
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
