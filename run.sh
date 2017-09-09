#!/bin/sh

go get -d -v
go build -buildmode=plugin -o hook/test.so hook/common.go hook/test.go
go build -buildmode=plugin -o hook/panic.so hook/common.go hook/panic.go
go run main.go
