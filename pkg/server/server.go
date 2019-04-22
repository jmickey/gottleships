package server

import (
	"bufio"
	"fmt"
	"net"
	"regexp"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/battleship"
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
			Trans: make(chan string),
			Recv:  make(chan string),
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
	logPrfx := fmt.Sprintf("[client %v]", c.Conn.RemoteAddr().String())
	defer logger.Infof("%v closed the connection", logPrfx)

	var gm *battleship.Game
	go receiver(c, logPrfx)
	go sender(c, logPrfx)

	select {
	case msg := <-c.Recv:
		switch msg {

		case "START GAME":
			gm = battleship.NewGame()
			c.Trans <- "POSITIONING SHIPS"
			c.Trans <- "SHIPS IN POSITION"

		default:
			valid, err := regexp.MatchString("^[A-I][1-9]$", msg)
			if err != nil || !valid {
				logger.Fatalf("%v received invalid msg: %v", logPrfx, msg)
			}

			if gm.Fire(msg) {

			}
		}
	}

}

func receiver(c *client.Client, logPrfx string) {
	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		msg := scanner.Text()
		logger.Infof("%v received message from client: %v", logPrfx, msg)
		c.Recv <- msg
	}
}

func sender(c *client.Client, logPrfx string) {
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
