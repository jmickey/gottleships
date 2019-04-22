package client

import (
	"fmt"
	"log"
	"net"

	"github.com/google/logger"
)

// Client represents a client connection
type Client struct {
	Conn  net.Conn
	Trans chan []byte
	Recv  chan []byte
}

// StartClient starts the application in client mode
func StartClient(hostname string, port string) error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%s", hostname, port))
	if err != nil {
		return fmt.Errorf("Couldn't connect to the server. %v", err)
	}

	client := &Client{
		Conn:  conn,
		Trans: make(chan []byte),
		Recv:  make(chan []byte),
	}

	go client.Receive()
	go client.Send()

	return nil
}

func (c *Client) Receive() {
	for {
		select {
		case msg, ok := <-c.Recv:
			if !ok {
				logger.Fatal("Channel closed")
			}
			fmt.Printf("Received message to send: %s", msg)
		}
	}
}

func (c *Client) Send() {
	for {
		select {
		case msg, ok := <-c.Trans:
			if !ok {
				log.Fatal("Channel closed")
			}
			fmt.Printf("Received message: %s", msg)
		}
	}
}
