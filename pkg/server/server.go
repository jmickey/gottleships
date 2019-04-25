package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/battleship"
	"github.com/jaymickey/gottleships/pkg/client"
)

// StartServer starts the application in server mode
func StartServer(port string) error {
	ln, err := net.Listen("tcp4", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		logger.Fatalf("%v", err.Error())
	}
	logger.Infof("Listening on localhost:%v", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		logger.Infof("Client connected, addr: %s", conn.RemoteAddr().String())
		c := &client.Client{
			Conn:   conn,
			Send:   make(chan string),
			Recv:   make(chan string),
			Closed: make(chan bool),
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

	var gm *battleship.Game
	go receiver(c, logPrfx)
	go sender(c, logPrfx)

	for {
		select {
		case <-c.Closed:
			close(c.Send)
			return

		case msg, ok := <-c.Recv:
			switch msg {

			case "":
				if !ok {
					close(c.Send)
					return
				}

			case "START GAME":
				gm = battleship.NewGame()
				c.Send <- "POSITIONING SHIPS"
				c.Send <- "SHIPS IN POSITION"

			default:
				valid, err := regexp.MatchString("^[A-I][1-9]$", msg)
				if err != nil || !valid {
					close(c.Send)
					logger.Errorf("%v received invalid msg: '%v', closing connection", logPrfx, msg)
					return
				}

				hit, err := gm.Fire(msg)
				if err != nil {
					close(c.Send)
					logger.Errorf("%v received invalid msg: '%v', closing connection", logPrfx, msg)
					return
				}

				switch hit {
				case true:
					c.Send <- "HIT"
					if gm.IsGameOver() {
						c.Send <- strconv.Itoa(gm.Shots())
						return
					}
				default:
					c.Send <- "MISS"
				}
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
				logger.Infof("%v connection closed by client", logPrfx)
				close(c.Recv)
				return
			}
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
		case msg, ok := <-c.Send:
			if !ok {
				return
			}
			logger.Infof("%v sending message: %s", logPrfx, msg)
			wr.WriteString(fmt.Sprintf("%s\n", msg))
			wr.Flush()
		}
	}
}
