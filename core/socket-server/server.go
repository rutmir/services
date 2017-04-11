package net_server

import (
	"errors"
	"fmt"
	"net"
)

type OnConnectionHandler func(*NetSocket)
type OnDataHandler func([]byte)
type OnErrorHandler func(error)
type OnCloseHandler func()
type ServerType string

// NetServer type for net server instance
type NetServer struct {
	Type       ServerType
	Host       string
	Port       int
	BufferSize int

	OnConnection OnConnectionHandler
}

// Run start listening on providen port
func (ns *NetServer) Run() error {
	err := validateNetServer(ns)
	if err != nil {
		return err
	}

	srv, err := net.Listen(string(ns.Type), fmt.Sprintf("%s:%v", ns.Host, ns.Port))
	if err != nil {
		return err
	}
	defer srv.Close()

	for {
		conn, err := srv.Accept()
		if err != nil {
			return err
		}

		sock := new(NetSocket)
		sock.connection = conn
		ns.OnConnection(sock)

		go handleConnection(sock, ns.BufferSize)
	}
}

// NetSocket type for net socket instance
type NetSocket struct {
	connection net.Conn
	OnData     OnDataHandler
	OnClose    OnCloseHandler
	OnError    OnErrorHandler
}

// Write bytes into socket
func (ns *NetSocket) Write(b []byte) (n int, err error) {
	return ns.connection.Write(b)
}

// Close socket
func (ns *NetSocket) Close() error {
	if ns.OnClose != nil {
		ns.OnClose()
	}

	return ns.connection.Close()
}

const (
	ServerTypeTcp ServerType = "tcp" // ServerTypeTcp strong typing for TCP server
	ServerTypeUdp ServerType = "udp" // ServerTypeUdp strong typing for UPD server
)

// CreateServer return new instance of NetServer
func CreateServer() *NetServer {
	return new(NetServer)
}

func handleConnection(socket *NetSocket, bufferSize int) {
	buf := make([]byte, bufferSize)
	defer socket.Close()

	for {
		reqLen, err := socket.connection.Read(buf)
		if err != nil {
			if socket.OnError != nil {
				socket.OnError(err)
			}
			break
		}

		if reqLen == 0 {
			break
		} else {
			socket.OnData(buf[:reqLen])
		}
		buf = make([]byte, bufferSize)
	}
}

func validateNetServer(ns *NetServer) error {
	if len(ns.Type) == 0 {
		return errors.New("Type is required")
	}
	if len(ns.Host) == 0 {
		return errors.New("Host is required")
	}
	if ns.Port < 1 || ns.Port > 65535 {
		return errors.New("Wrong Port value")
	}
	if ns.BufferSize < 256 {
		return errors.New("Too small BufferSize")
	}
	if ns.OnConnection == nil {
		return errors.New("OnConnection handler required")
	}
	return nil
}
