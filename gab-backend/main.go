package main

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/rutmir/services/core/log"
	ns "github.com/rutmir/services/core/net-server"
	dto "github.com/rutmir/services/entities/dto/v2"
)

var server *ns.NetServer

func main() {
	log.Info("Initialize Gab End point")
	server = ns.CreateServer()
	server.Type = ns.ServerTypeTcp
	server.Host = "localhost"
	server.Port = 3333
	server.BufferSize = 4096

	server.OnConnection = func(socket *ns.NetSocket) {
		ep := new(Endpoint)

		ep.Socket = socket

		socket.OnError = func(err error) {
			log.Err("sock 'on error': %v", err)
		}

		socket.OnClose = func(){
			log.Info("Socket closed")
		}

		socket.OnData = func(data []byte) {
			source := append(ep.Chunk, data...)
			hasData := true

			for hasData == true {
				size := int16(binary.LittleEndian.Uint16(source[0:2])) // read first two bites
				length := len(source)

				log.Info("size: %v", size)

				if size > 0 && length >= int(size + 2) {
					log.Info("Extract package")
					pack := new(dto.Package)

					if err := proto.Unmarshal(source[2:size + 2], pack); err != nil {
						log.Warn(err)
					} else {
						ep.HandlePackage(pack)
					}

					source = source[size + 2:]
				} else {
					log.Info("Save chunk")
					ep.Chunk = source;
					hasData = false;
				}
			}
		}
	}

	log.Info("Start Gab End point")
	err := server.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

type Endpoint struct {
	Socket *ns.NetSocket
	Chunk  []byte
}

// Process single package
func (ep *Endpoint) HandlePackage(pack *dto.Package) bool {
	log.Info(pack.ToFullString())
	return false;
}
