package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/client"
)

// StartServer starts the application in server mode
func StartServer(port string) error {
	ln, _ := net.Listen("tcp4", fmt.Sprintf("localhost:%v", port))
	logger.Infof("Listening on localhost:%v", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		logger.Infof("Client connected, addr: %s", conn.RemoteAddr().String())
		c := &client.Client{
			Conn:  conn,
			Trans: make(chan []byte),
			Recv:  make(chan []byte),
		}
		// Launch goroutine to handle connection. Allowing the server, freeing
		// the server to accept another connection.
		go connHandler(c)
	}
}

// Function to handle connections as they are received. Allows for concurrency when
// called as a goroutine. In the event of an error the error message is logged and the
// function is returned. When the function is returned the connection is automatically
// closed due to the defer statement. The next time the client attempts to read or
// write to the connection they will receive an io.EOF error, signalling the connection
// is closed.
func connHandler(c *client.Client) {
	defer c.Conn.Close()
	logPrfx := fmt.Sprintf("client %v:", c.Conn.RemoteAddr().String())
	defer logger.Infof("%v closed the connection", logPrfx)

	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		msg := scanner.Text()
		logger.Infof("%v received message '%v'", logPrfx, msg)
	}
}
