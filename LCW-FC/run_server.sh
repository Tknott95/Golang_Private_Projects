#!/bin/bash

# Assuming you already got project
cd $GOPATH/src/github.com/tknott95/LCW-FC/

go fmt .

go build

go run application.go