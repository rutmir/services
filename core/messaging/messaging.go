package messaging

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/rutmir/services/core/log"
	"github.com/streadway/amqp"
)

// EventBusInterface is used for all AMQP interaction
type EventBusInterface struct {
	Available bool
	Conn      *amqp.Connection

	host     string
	port     string
	user     string
	password string
	channel  *amqp.Channel
	queue    amqp.Queue
	listener func(map[string]interface{}, string)
}

// GetInstance - Initialize EventBusInterface, connects to RabbitMQ
func GetInstance() *EventBusInterface {
	log.Info("Initializing Event Bus...")
	eventBus := new(EventBusInterface)

	eventBus.Available = false

	eventBus.port = os.Getenv("AMQP_PORT")
	if len(eventBus.port) == 0 {
		if eventBus.failOnError(errors.New("initialization error"), "Required to set 'AMQP_PORT' environment") {
			return eventBus
		}
	}

	eventBus.host = os.Getenv("AMQP_HOST")
	if len(eventBus.host) == 0 {
		eventBus.host = os.Getenv("HOSTNAME")
	}
	if len(eventBus.host) == 0 {
		if eventBus.failOnError(errors.New("initialization error"), "Required to set 'AMQP_HOST' environment") {
			return eventBus
		}
	}

	eventBus.user = os.Getenv("AMQP_USERNAME")
	if len(eventBus.user) == 0 {
		if eventBus.failOnError(errors.New("initialization error"), "Required to set 'AMQP_USERNAME' environment") {
			return eventBus
		}
	}

	eventBus.password = os.Getenv("AMQP_PASS")
	if len(eventBus.password) == 0 {
		if eventBus.failOnError(errors.New("initialization error"), "Required to set 'AMQP_PASS' environment") {
			return eventBus
		}
	}

	//Make connection
	amqpServerURL := "amqp://" + eventBus.user + ":" + eventBus.password + "@" + eventBus.host + ":" + eventBus.port + "/"
	log.Info(amqpServerURL)
	conn, err := amqp.Dial(amqpServerURL)
	if eventBus.failOnError(err, "failed to connect to RabbitMQ") {
		return eventBus
	}
	eventBus.Conn = conn
	//defer conn.Close()

	ch, err := conn.Channel()
	if eventBus.failOnError(err, "failed to open a channel") {
		return eventBus
	}
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		"optima.eventbus", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if eventBus.failOnError(err, "failed to declare an exchange") {
		return eventBus
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if eventBus.failOnError(err, "failed to declare a queue") {
		return eventBus
	}

	eventBus.channel = ch
	eventBus.queue = q
	//	eventBus.listener = make(map[string]func(map[string]interface{}, string))
	go eventBus.startListening()
	eventBus.Available = true

	return eventBus
}

// BroadcastToAll sends message to all listeners
func (eventBus *EventBusInterface) BroadcastToAll(message map[string]interface{}) (bool, error) {
	//optima.eventbus.common
	message["__group"] = "all"
	return eventBus.sendMessage(message, "optima.eventbus", "optima.eventbus.common")
}

// BroadcastInternal sends message internal channel
func (eventBus *EventBusInterface) BroadcastInternal(message map[string]interface{}) (bool, error) {
	//optima.eventbus.private
	message["__group"] = "internal"
	return eventBus.sendMessage(message, "optima.eventbus", "optima.eventbus.private")
}

// BroadcastToSubscribers sends message to a group specified by routingKey
func (eventBus *EventBusInterface) BroadcastToSubscribers(message map[string]interface{}, routingKey string) (bool, error) {
	message["__group"] = routingKey
	return eventBus.sendMessage(message, "optima.eventbus", routingKey)
}

func (eventBus *EventBusInterface) sendMessage(message map[string]interface{}, exchange string, routingKey string) (bool, error) {
	if !eventBus.Available {
		return false, nil
	}

	body, err := json.Marshal(message)
	log.Info("Sending Event: %v %v %v", exchange, routingKey, string(body))
	if err != nil {
		return false, err
	}
	err = eventBus.channel.Publish(
		exchange,
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return false, err
	}
	return true, nil
}

// AttachListener registers event listener
func (eventBus *EventBusInterface) AttachListener(listener func(map[string]interface{}, string)) {
	eventBus.listener = listener
}

/*
func (eventBus *EventBusInterface) RemoveListener(id string) {
	delete(eventBus.listeners, id)
}
*/

func (eventBus *EventBusInterface) startListening() {
	err := eventBus.channel.QueueBind(
		eventBus.queue.Name,       // queue name
		"optima.eventbus.private", // routing key
		"optima.eventbus",         // exchange
		false,
		nil)
	eventBus.failOnError(err, "Failed to bind a queue")

	err = eventBus.channel.QueueBind(
		eventBus.queue.Name,      // queue name
		"optima.eventbus.common", // routing key
		"optima.eventbus",        // exchange
		false,
		nil)
	eventBus.failOnError(err, "failed to bind a queue")

	msgs, err := eventBus.channel.Consume(
		eventBus.queue.Name, // queue
		"",                  // consumer
		true,                // auto ack
		false,               // exclusive
		false,               // no local
		false,               // no wait
		nil,                 // args
	)
	eventBus.failOnError(err, "failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if eventBus.listener == nil {
				log.Info("Message received but no listener attached.")
				continue
			}
			log.Info("Sending event to listener...")
			var message map[string]interface{}
			message = make(map[string]interface{})
			json.Unmarshal([]byte(d.Body), &message)
			if str, ok := message["__group"].(string); ok {
				eventBus.listener(message, str)
			} else {
				eventBus.listener(message, "n/a")
			}
		}
	}()

	log.Info("[*]Listening for messages.")
	<-forever
}

// SubscribeListener registers queues launches consumer loop
func (eventBus *EventBusInterface) SubscribeListener(routingKey string, listener func(message map[string]interface{}, group string)) {
	if !eventBus.Available {
		return
	}
	ch, err := eventBus.Conn.Channel()
	if eventBus.failOnError(err, "failed to open a channel") {
		return
	}
	//defer ch.Close()

	err = ch.ExchangeDeclare(
		"optima.eventbus", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if eventBus.failOnError(err, "failed to declare an exchange") {
		return
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if eventBus.failOnError(err, "failed to declare a queue") {
		return
	}

	err = ch.QueueBind(
		q.Name,            // queue name
		routingKey,        // routing key
		"optima.eventbus", // exchange
		false,
		nil)
	if eventBus.failOnError(err, "failed to bind a queue") {
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if eventBus.failOnError(err, "failed to register a consumer") {
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Info("Sending event to Group Listener: %s", routingKey)
			var message map[string]interface{}
			message = make(map[string]interface{})
			json.Unmarshal([]byte(d.Body), &message)
			if str, ok := message["__group"].(string); ok {
				listener(message, str)
			} else {
				listener(message, "n/a")
			}
		}
	}()

	log.Info("[%s] Listening for messages.", routingKey)
	<-forever
}

func (eventBus *EventBusInterface) failOnError(err error, msg string) bool {
	if err != nil {
		log.Fatal("%s: %v", msg, err)
		return true
	}
	return false
}
