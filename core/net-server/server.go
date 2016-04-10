package net_server

import (
	"fmt"
	"net"
)

type OnConnectionHandler func(*NetSocket)
type OnDataHandler func([]byte)
type OnErrorHandler func(error)
type OnCloseHandler func()
type ServerType string

type NetServer struct {
	Type         ServerType
	Host         string
	Port         int
	BufferSize   int

	OnConnection OnConnectionHandler
}

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

type NetSocket struct {
	connection net.Conn
	OnData     OnDataHandler
	OnClose    OnCloseHandler
	OnError    OnErrorHandler
}

func (ns *NetSocket) Write(b []byte) (n int, err error) {
	return ns.connection.Write(b)
}

func (ns *NetSocket) Close() error {
	if ns.OnClose != nil {
		ns.OnClose()
	}

	return ns.connection.Close()
}

const (
	ServerTypeTcp ServerType = "tcp"
	ServerTypeUdp ServerType = "udp"
)

func CreateServer() *NetServer {
	return new(NetServer)
}

/*func startServer(serverType ServerType, host, port string) {
	// Listen for incoming connections.
	l, err := net.Listen(string(serverType), host + ":" + port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + host + ":" + port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}*/

// Handles incoming requests.
func handleConnection(conn *NetSocket, bufferSize int) {
	buf := make([]byte, bufferSize)
	reqLen, err := conn.connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	conn.OnData(buf)

	fmt.Println("Read from socket:", reqLen)
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}

func validateNetServer(ns *NetServer) error {
	if len(ns.Type) == 0 {
		return fmt.Errorf("Type is required")
	}
	if len(ns.Host) == 0 {
		return fmt.Errorf("Host is required")
	}
	if ns.Port < 1 || ns.Port > 65535 {
		return fmt.Errorf("Wrong Port value")
	}
	if ns.BufferSize < 256 {
		return fmt.Errorf("Too small BufferSize")
	}
	if ns.OnConnection == nil {
		return fmt.Errorf("OnConnection handler required")
	}
	return nil
}