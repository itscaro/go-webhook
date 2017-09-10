#!/bin/sh

go get -d -v
go build -a -installsuffix cgo -buildmode=plugin -o hook/test.so hook/common.go hook/test.go
go build -a -installsuffix cgo -buildmode=plugin -o hook/panic.so hook/common.go hook/panic.go
go run main.go
