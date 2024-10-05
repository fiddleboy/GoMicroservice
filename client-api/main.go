package main

import (
	"context"
	"fmt"

	"github.com/fiddleboy/GoMicroservice/gRPC-play/protos/currency"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	currencyClient := currency.NewCurrencyClient(conn)
	rateReqeust := &currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value["EUR"]),
		Destination: currency.Currencies(currency.Currencies_value["USD"]),
	}
	resp, err := currencyClient.GetRate(context.Background(), rateReqeust)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The server returned rate: %f", resp.Rate)
}
