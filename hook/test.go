package main

import "C"

type testType struct{}

func (h testType) Exec(data []byte) interface{} {
	return JsonResponse{
		Message: "I received this data: " + string(data),
	}
}

var Test testType
