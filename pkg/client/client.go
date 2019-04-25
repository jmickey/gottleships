package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/ui"
)

// Client represents a client connection
type Client struct {
	Conn   net.Conn
	Send   chan string
	Recv   chan string
	Closed chan bool
}

// StartClient starts the application in client mode
func StartClient(hostname string, port string) error {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%s", hostname, port))
	if err != nil {
		return fmt.Errorf("Couldn't connect to the server. %v", err)
	}

	client := &Client{
		Conn:   conn,
		Send:   make(chan string),
		Recv:   make(chan string),
		Closed: make(chan bool),
	}

	// Launch goroutines to hand sending/receiving over the network
	go client.receive()
	go client.send()
	client.Send <- "START GAME"
	if err := ui.Load(client.Send, client.Recv, client.Closed); err != nil {
		return err
	}
	return nil
}

func (c *Client) receive() {
	rd := bufio.NewReader(c.Conn)
	for {
		msg, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				c.Closed <- true
				return
			}
			logger.Errorf("%v", err.Error())
			return
		}
		msg = strings.TrimSuffix(msg, "\n")
		logger.Infof("received message from server: %v", msg)
		if msg == "POSITIONING SHIPS" || msg == "SHIPS IN POSITION" {
			continue
		}
		c.Recv <- msg
	}
}

func (c *Client) send() {
	wr := bufio.NewWriter(c.Conn)
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return
			}
			logger.Infof("sending message: %s", msg)
			wr.WriteString(fmt.Sprintf("%s\n", msg))
			wr.Flush()
		}
	}
}
