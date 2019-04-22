// Copyright (C) 2019 Josh Michielsen <git@mickey.dev>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/logger"
	"github.com/jaymickey/gottleships/pkg/client"
	"github.com/jaymickey/gottleships/pkg/server"
	flag "github.com/spf13/pflag"
)

var (
	mode     string
	hostname string
	port     int
	logFile  string
)

func init() {
	flag.StringVarP(&mode, "mode", "m", "", "(Required) Application mode to launch. Valid input: 'server' or 'client'")
	flag.StringVarP(&hostname, "hostname", "h", "localhost", "(Optional) The hostname to listen via or connect")
	flag.IntVarP(&port, "port", "p", 8080, "(Optional) The port number to listen via or connect")
	flag.StringVarP(&logFile, "log", "l", "", "(Optional) File to direct log output (default \"stdout\")")
	flag.CommandLine.SortFlags = false
}

func main() {
	flag.Parse()

	// Setup log file if provided
	lf := ioutil.Discard
	if logFile != "" {
		var err error
		lf, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			logger.Fatalf("Failed to open log file: %v", err)
		}
	}
	log := logger.Init("BattleshipLogger", true, false, lf)
	defer log.Close()

	p := fmt.Sprintf("%v", port)

	switch strings.ToLower(mode) {
	case "client":
		log.Info("Starting application in client mode")
		if err := client.StartClient(hostname, p); err != nil {
			log.Fatal(err.Error())
		}
	case "server":
		log.Info("Starting application in server mode")
		server.StartServer(p)
	default:
		logger.Error("Invalid value for flag 'mode'")
		flag.Usage()
	}
}
