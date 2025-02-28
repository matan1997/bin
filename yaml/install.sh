#! /bin/bash

#build for all linux kernel
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o convert yaml/main.go     