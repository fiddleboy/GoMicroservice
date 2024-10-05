package main

import (
	"context"

	"github.com/fiddleboy/GoMicroservice/gRPC-play/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func subscribe(logger hclog.Logger, client currency.CurrencyClient) {
	sub, err := client.SubscribeRates(context.Background())
	if err != nil {
		logger.Error("Unable to call SubscribeRates(); connection cannot be established!")
	}
	for {
		response, err := sub.Recv()
		if response.GetError()
	}
}

func main() {
	logger := hclog.Default()
	conn, err := grpc.NewClient("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	currencyClient := currency.NewCurrencyClient(conn)

	subscribe(logger, currencyClient)

	// rateReqeust := &currency.RateRequest{
	// 	Base:        currency.Currencies(currency.Currencies_value["EUR"]),
	// 	Destination: currency.Currencies(currency.Currencies_value["EUR"]),
	// }
	// resp, err := currencyClient.GetRate(context.Background(), rateReqeust)
	// if err != nil {
	// 	if s, ok := status.FromError(err); ok {
	// 		metadata := s.Details()[0].(*currency.RateRequest)
	// 		if s.Code() == codes.InvalidArgument {
	// 			fmt.Printf("Provided argument is invalid. with base: %v; dest: %v\n", metadata.Base, metadata.Base)
	// 			fmt.Println(*metadata)
	// 			return
	// 		}
	// 	}
	// 	return
	// }
	// fmt.Println(resp.String())
	// fmt.Printf("The server returned rate: %f", resp.Rate)
}
