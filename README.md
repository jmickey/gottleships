# Gottleships: Go(Ba)ttleships

Single player implementation of the Battleships game in Go with a client/server architecture that communicates over raw sockets, complete with a terminal based UI (TUI). I originally created Gottleships as part of an assignment for university (UNE COSC340).

## Building & Running

### Prerequisites

- The [`Go` programming language](https://golang.org/dl/)
- `$GOPATH` configured (generally `~/go/`)
- Access to the internet.
  - See [here](https://github.com/golang/go/wiki/GoGetProxyConfig) for information regarding configuring a proxy if required.

### Build the Application

1. Clone the repo. There are 2 ways to do this:
    - `go get github.com/jaymickey/gottleships`
    - `git clone (https://)(git@)github.com:jaymickey/gottleships.git $GOPATH/src/github.com/jaymickey/gottleships`
2. `cd $GOPATH/src/github.com/jaymickey/gottleships`
3. `go install`

### Running Gottleships

Go install will download any dependencies, build the binary, and copy it into `$GOPATH/bin`. If `$GOPATH/bin` is in your `$PATH` variable, then you should be able to simple run `gottleships`, along with the required flags.

#### Flags

| Flag      | Short | Required? | Default     | Description                                                   |
|-----------|-------|-----------|-------------|---------------------------------------------------------------|
| -mode     | -m    | Yes       |             | Application mode to launch. Valid input: `server` or `client` |
| -hostname | -h    | No        | `localhost` | The hostname to listen via or connect                         |
| -port     | -p    | No        | 8080        | The port number to listen or connect                          |
| -log      | -l    | No        | `Stdout`    | Location of file to write log output                          |
