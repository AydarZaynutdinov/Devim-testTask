package main

import (
	"net/http"
	"os"
	"testTask/handler"
)

const (
	defaultAddress = "localhost:8080"
	addressKey     = "ADDRESS"
)

func main() {
	address := os.Getenv(addressKey)
	if address == "" {
		address = defaultAddress
	}

	http.HandleFunc("/", handler.GetValuteValueByCoordinate)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
