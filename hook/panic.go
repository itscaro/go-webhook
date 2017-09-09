package main

import "C"

type panicType struct{}


func (h panicType) Exec(data []byte) interface{} {
	panic("Panic simulation")
}

var Panic panicType
