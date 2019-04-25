#!/bin/bash
set -e

if [ "$#" -ne 2 ]; then
    echo "Error: Invalid number of args. Expected 2, got: $#"
    echo "Usage: startClient.sh [hostname] [port]"
    exit 1
fi

if ! [[ $2 =~ ^-?[0-9]+$ ]]; then
    echo "Error: $1 is not a valid port number."
    exit 1
fi

if [[ ! -x $(which go) ]]; then
    echo "Go not found, please install or ensure it is in your path"
    exit 1
fi

if [[ ! -d $GOPATH/src/github.com/jaymickey/gottleships ]]; then
    go get -u github.com/jaymickey/gottleships
fi

if [[ -x $(which gottleships) ]]; then
    gottleships -m client -h $1 -p $2
    else
        if [[ -x $(which $(go env GOPATH)/bin/gottleships) ]]; then
            $(go env GOPATH)/bin/gottleships -m client -h $1 -p $2
        else
            echo "Couldn't find binary for Gottleships, check the README"
            exit 1
        fi
    fi
fi