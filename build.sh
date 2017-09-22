#!/bin/bash

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.build=$DRONE_BUILD_NUMBER" -a -tags netgo
