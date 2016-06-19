package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/rutmir/services/core/log"
	"github.com/rutmir/services/core/memcache"
	ns "github.com/rutmir/services/core/net-server"
	dto "github.com/rutmir/services/entities/dto/v2"
)

const (
	atPrefix = "at_"
	rtPrefix = "rt_"
)

var memCtrl memcache.MemCache
var server *ns.NetServer

func main() {
	log.Info("Initialize Chat End point")

	var err error
	memCtrl, err = memcache.GetLocalInstance("memcached", "test")
	if err != nil {
		log.Fatal(err)
	}

	server = ns.CreateServer()
	server.Type = ns.ServerTypeTcp
	server.Host = "localhost"
	server.Port = 3333
	server.BufferSize = 4096

	server.OnConnection = func(socket *ns.NetSocket) {
		ep := new(Endpoint)

		ep.Socket = socket
		ep.packageProcessor = ep.handleUnauthorizedPackage
		ep.BufferSize = 4096
		ep.msgHeadersMap = make(map[string]*dto.Header)
		ep.msgBodyMap = make(map[string]map[int][]byte)

		if err := ep.initAmqp(); err != nil {

		}

		socket.OnError = func(err error) {
			log.Err("sock 'on error': %v", err)
		}

		socket.OnClose = func() {
			log.Info("Socket closed")
			ep.dispose()
		}

		socket.OnData = func(data []byte) {
			source := append(ep.chunk, data...)
			hasData := true

			for hasData == true {
				size := int16(binary.LittleEndian.Uint16(source[0:2])) // read first two bites
				length := len(source)

				log.Info("size: %v", size)

				if size > 0 && length >= int(size+2) {
					log.Info("Extract package")
					pack := new(dto.Package)

					if err := proto.Unmarshal(source[2:size+2], pack); err != nil {
						log.Warn(err)
					} else {
						ep.packageProcessor(pack)
					}

					source = source[size+2:]
				} else {
					log.Info("Save chunk")
					ep.chunk = source
					hasData = false
				}
			}
		}
	}

	log.Info("Start Chat End point")
	if err := server.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

type funcPackageHandle func(pack *dto.Package)

// Endpoint - realisation of chat socket endpoint
type Endpoint struct {
	Socket           *ns.NetSocket
	BufferSize       int
	Auth             *dto.AuthTokenMem
	chunk            []byte
	packagePipe      chan *dto.Package
	packageProcessor funcPackageHandle
	msgHeadersMap    map[string]*dto.Header
	msgBodyMap       map[string]map[int][]byte
	amqpConnection   *amqp.Connection
	amqpChannel      *amqp.Channel
	amqpQueue        amqp.Queue
	amqpMessages     <-chan amqp.Delivery
}

func (ep *Endpoint) initAmqp() error {
	host := os.Getenv("AMQP_HOST")
	if len(host) == 0 {
		return fmt.Errorf("AMQP error: Required to set 'AMQP_HOST' environment")
	}

	port := os.Getenv("AMQP_PORT")
	if len(port) == 0 {
		return fmt.Errorf("AMQP error: Required to set 'AMQP_PORT' environment")
	}

	user := os.Getenv("AMQP_USERNAME")
	if len(user) == 0 {
		return fmt.Errorf("AMQP error: Required to set 'AMQP_USERNAME' environment")
	}

	password := os.Getenv("AMQP_PASS")
	if len(password) == 0 {
		return fmt.Errorf("AMQP error: Required to set 'AMQP_PASS' environment")
	}

	amqpServerURL := "amqp://" + user + ":" + password + "@" + host + ":" + port + "/"

	log.Info(amqpServerURL)

	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return err // "Failed to connect to RabbitMQ"
	}

	ch, err := conn.Channel()
	if err != nil {
		defer conn.Close()
		return err // "Failed to open a channel"
	}

	ep.amqpConnection = conn
	ep.amqpChannel = ch

	return nil
}

func (ep *Endpoint) handleUnauthorizedPackage(pack *dto.Package) {
	log.Info(pack.ToString())

	if len(pack.MessageID) == 0 || !pack.Closed {
		ep.responseBadRequest(pack.MessageID)
		return
	}

	head := new(dto.Header)
	if err := proto.Unmarshal(pack.Data, head); err != nil {
		log.Err(err)
		ep.responseBadRequest(pack.MessageID)
		return
	}
	if head.Action != dto.Action_Authorize || len(head.Meta) == 0 {
		ep.responseUnauthorized(pack.MessageID)
		ep.destroyEndPoint()
		return
	}

	mem, err := memCtrl.Get(atPrefix + head.Meta)
	if err != nil {
		log.Err(err)
		ep.responseUnauthorized(pack.MessageID)
		ep.destroyEndPoint()
		return
	}

	auth := new(dto.AuthTokenMem)
	if err := proto.Unmarshal(mem.Value, auth); err != nil {
		log.Err(err)
		ep.responseUnauthorized(pack.MessageID)
		ep.destroyEndPoint()
		return
	}
	ep.Auth = auth
	q, err := ep.amqpChannel.QueueDeclare(
		ep.Auth.AccessToken, // name
		false,               // durable
		false,               // autoDelete
		true,                // exclusive
		false,               // noWait
		nil,                 // arguments
	)

	if err != nil {
		log.Err(err) // "Failed to declare a queue"
		ep.responseInternalError(pack.MessageID)
		ep.destroyEndPoint()
		return
	}

	ep.amqpQueue = q
	msgs, err := ep.amqpChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Err(err) // "Failed to register a consumer"
		ep.responseInternalError(pack.MessageID)
		ep.destroyEndPoint()
		return
	}

	ep.amqpMessages = msgs
	ep.packagePipe = make(chan *dto.Package)

	go ep.runWorkerPipe()
	go ep.runPackagePipe()
	ep.packageProcessor = ep.handlePackage

	log.Info(head.ToString())
}

func (ep *Endpoint) handlePackage(pack *dto.Package) {
	log.Info(pack.ToString())

	ep.packagePipe <- pack
}

func (ep *Endpoint) handleResponse(messageID string, head *dto.Header, body []byte) {
	if h, err := proto.Marshal(head); err != nil {
		log.Warn(err)
	} else {
		if body != nil && len(body) > 0 {
			pack := new(dto.Package)
			pack.MessageID = messageID
			pack.PackageNo = 0
			pack.Closed = false
			pack.Data = h

			if p, err := proto.Marshal(pack); err != nil {
				log.Err(err)
			} else {
				ep.writeToSocket(p)
				lastSentIdx := 0
				length := len(body)

				for lastSentIdx < length {
					nextSentCnt := ep.BufferSize - 2
					if lastSentIdx+nextSentCnt >= length {
						nextSentCnt = length - lastSentIdx
						pack.Closed = true
					}
					pack.PackageNo++
					pack.Data = body[lastSentIdx : lastSentIdx+nextSentCnt]

					if p, err := proto.Marshal(pack); err != nil {
						log.Err(err)
					} else {
						ep.writeToSocket(p)
					}

					lastSentIdx += nextSentCnt
				}
			}
		} else {
			pack := new(dto.Package)
			pack.MessageID = messageID
			pack.PackageNo = 0
			pack.Closed = true
			pack.Data = h

			if p, err := proto.Marshal(pack); err != nil {
				log.Err(err)
			} else {
				ep.writeToSocket(p)
			}
		}
	}
}

func (ep *Endpoint) responseUnauthorized(messageID string) {
	log.Info("Unauthorized client response")
	head := new(dto.Header)
	head.Action = dto.Action_Result
	head.Timestamp = time.Now().UnixNano()

	result := new(dto.Result)
	result.Code = 401
	result.Result = "Error"
	result.Message = "Unauthorized"

	if data, err := proto.Marshal(result); err != nil {
		log.Warn(err)
	} else {
		ep.handleResponse(messageID, head, data)
	}
}

func (ep *Endpoint) responseBadRequest(messageID string) {
	log.Info("Bad request response")
	head := new(dto.Header)
	head.Action = dto.Action_Result
	head.Timestamp = time.Now().UnixNano()

	result := new(dto.Result)
	result.Code = 400
	result.Result = "Error"
	result.Message = "Bad Request"

	if data, err := proto.Marshal(result); err != nil {
		log.Err(err)
	} else {
		ep.handleResponse(messageID, head, data)
	}
}

func (ep *Endpoint) responseInternalError(messageID string) {
	log.Info("Internal server response")
	head := new(dto.Header)
	head.Action = dto.Action_Result
	head.Timestamp = time.Now().UnixNano()

	result := new(dto.Result)
	result.Code = 500
	result.Result = "Error"
	result.Message = "Internal server error"

	if data, err := proto.Marshal(result); err != nil {
		log.Err(err)
	} else {
		ep.handleResponse(messageID, head, data)
	}
}

func (ep *Endpoint) writeToSocket(data []byte) {
	//log.Info("len: %v, data: %v", len(data), data)
	pref := []byte{0, 0}
	binary.LittleEndian.PutUint16(pref, uint16(len(data)))
	//log.Info("ulen: %v, pref: %v", uint16(len(data)), pref)
	if _, err := ep.Socket.Write(append(pref, data...)); err != nil {
		log.Err(err)
	}
}

func (ep *Endpoint) processMessage(messageID string) {
	log.Info("messageID: %v", messageID)

	h := ep.msgHeadersMap[messageID]
	rawBody := ep.msgBodyMap[messageID]

	delete(ep.msgHeadersMap, messageID)
	delete(ep.msgBodyMap, messageID)

	if h != nil {
		msg := new(dto.InternalMessage)
		msg.Header = h

		if rawBody != nil {
			var b []byte
			var keys []int

			for k := range rawBody {
				keys = append(keys, k)
			}
			sort.Ints(keys)

			for _, k := range keys {
				b = append(b, rawBody[k]...)
			}
			msg.Body = b
		}

		if data, err := proto.Marshal(msg); err != nil {
			log.Err(err)
			//TODO return error response
		} else {
			err := ep.amqpChannel.Publish(
				"",            // exchange
				"chat_worker", // routing key
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType:   "application/x-protobuf",
					CorrelationId: messageID,
					ReplyTo:       ep.amqpQueue.Name,
					Body:          data,
				})
			if err != nil {
				log.Err(err)
				//TODO return error response
			}
		}
	}
}

func (ep *Endpoint) runPackagePipe() {
	for pack := range ep.packagePipe {
		if len(pack.MessageID) == 0 {
			ep.responseBadRequest(pack.MessageID)
			continue
		}

		if pack.PackageNo == 0 {
			head := new(dto.Header)
			if err := proto.Unmarshal(pack.Data, head); err != nil {
				log.Err(err)
				ep.responseBadRequest(pack.MessageID)
				continue
			}
			ep.msgHeadersMap[pack.MessageID] = head
		} else {
			bodyChain := ep.msgBodyMap[pack.MessageID]
			if bodyChain == nil {
				bodyChain = make(map[int][]byte)
				ep.msgBodyMap[pack.MessageID] = bodyChain
			}
			bodyChain[int(pack.PackageNo)] = pack.Data
		}

		if pack.Closed {
			ep.processMessage(pack.MessageID)
		}
	}
}

func (ep *Endpoint) runWorkerPipe() {
	for d := range ep.amqpMessages {
		log.Info("Response form worker with: %s", d.CorrelationId)
	}
}

func (ep *Endpoint) destroyEndPoint() {
	log.Info("Destroy socked")
	ep.Socket.Close()
}

func (ep *Endpoint) dispose() {
	if ep.amqpChannel != nil {
		ep.amqpChannel.Close()
	}

	if ep.amqpConnection != nil {
		ep.amqpConnection.Close()
	}

	if ep.packagePipe != nil {
		close(ep.packagePipe)
	}

	//if ep.amqpMessages != nil {
	//	close(ep.amqpMessages)
	//}
}
