package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/client"
	"github.com/jaymickey/gottleships/pkg/server"
)

var (
	mode    = flag.String("mode", "", "Required. Application mode to launch. Valid input: 'server' or 'client'.")
	logFile = flag.String("log", "stdout", "Optional. Log to a file.")
)

func main() {
	flag.Parse()

	// Setup log file if provided
	lf := ioutil.Discard
	if *logFile != "" {
		var err error
		lf, err = os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			logger.Fatalf("Failed to open log file: %v", err)
		}
	}
	log := logger.Init("BattleshipLogger", true, false, lf)
	defer log.Close()
	log.Infof("Starting in %s mode", *mode)

	switch *mode {
	case "client":
		if err := client.StartClient("localhost", "8080"); err != nil {
			log.Fatal(err.Error())
		}
	case "server":
		server.StartServer("8080")
	}
}
