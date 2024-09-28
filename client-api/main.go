package main

import (
	"gRPC-play/protos/currency"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:9092")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	currencyClient := currency.New

}
