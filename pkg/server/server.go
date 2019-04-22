package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"

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
	case msg, ok := <-c.Recv:
		switch msg {

		case "":
			if !ok {
				return
			}

		case "START GAME":
			gm = battleship.NewGame()
			c.Trans <- "POSITIONING SHIPS"
			c.Trans <- "SHIPS IN POSITION"

		default:
			valid, err := regexp.MatchString("^[A-I][1-9]$", msg)
			if err != nil || !valid {
				close(c.Trans)
				close(c.Recv)
				logger.Errorf("%v received invalid msg: '%v', closing connection", logPrfx, msg)
				return
			}

			switch gm.Fire(msg) {

			case true:
				c.Trans <- "HIT"
				if gm.IsGameOver() {
					c.Trans <- string(gm.Shots())
					return
				}

			default:
				c.Trans <- "MISS"
			}
		}
	}
}

func receiver(c *client.Client, logPrfx string) {
	rd := bufio.NewReader(c.Conn)
	for {
		msg, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Errorf("%v connection closed by client", logPrfx)
				return
			}
			logger.Errorf("%v %v", logPrfx, err.Error())
			return
		}
		msg = strings.TrimSuffix(msg, "\n")
		logger.Infof("%v received message from client: %v", logPrfx, msg)
		c.Recv <- msg
	}
}

func sender(c *client.Client, logPrfx string) {
	wr := bufio.NewWriter(c.Conn)
	for {
		select {
		case msg, ok := <-c.Trans:
			if !ok {
				logger.Errorf("%v transmit channel closed", logPrfx)
				return
			}
			logger.Infof("%v sending message: %s", logPrfx, msg)
			wr.WriteString(fmt.Sprintf("%s\n", msg))
			wr.Flush()
		}
	}
}
