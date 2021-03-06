package net_server

import (
	"fmt"
	"testing"
)

var server *NetServer

func TestMain(m *testing.M) {

	server = CreateServer()
	server.Type = ServerTypeTcp
	server.Host = "localhost"
	server.Port = 3333
	server.BufferSize = 1024

	server.OnConnection = func(socket *NetSocket) {
		socket.OnError = func(err error) {
			fmt.Println(err.Error())
		}

		socket.OnData = func(data []byte) {
			fmt.Println("I get new data.")
			socket.Write([]byte("Message received."))
		}
	}

	err := server.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestInitial(t *testing.T) {
	if server == nil {
		t.Fatal("failed to Create")
	}
}
