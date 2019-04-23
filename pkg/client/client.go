package client

import (
	"fmt"
	"net"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/ui"
)

// Client represents a client connection
type Client struct {
	Conn  net.Conn
	Trans chan string
	Recv  chan string
}

// StartClient starts the application in client mode
func StartClient(hostname string, port string) error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%s", hostname, port))
	if err != nil {
		return fmt.Errorf("Couldn't connect to the server. %v", err)
	}

	client := &Client{
		Conn:  conn,
		Trans: make(chan string),
		Recv:  make(chan string),
	}

	go client.receive()
	go client.send()
	if err = ui.Load(); err != nil {
		return err
	}

	return nil
}

func (c *Client) receive() {
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

func (c *Client) send() {
	for {
		select {
		case msg, ok := <-c.Trans:
			if !ok {
				logger.Fatal("Channel closed")
			}
			fmt.Printf("Received message: %s", msg)
		}
	}
}
