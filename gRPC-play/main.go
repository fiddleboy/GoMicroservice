package main

import (
	"net"
	"os"

	"github.com/fiddleboy/GoMicroservice/gRPC-play/data"
	"github.com/fiddleboy/GoMicroservice/gRPC-play/protos/currency"
	"github.com/fiddleboy/GoMicroservice/gRPC-play/server"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	gs := grpc.NewServer()
	er, err := data.NewRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}
	cs := server.NewCurrency(log, er)
	currency.RegisterCurrencyServer(gs, cs)
	reflection.Register(gs)

	// set up a tcp connection
	listener, err := net.Listen("tcp", ":9092")

	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	// have the gRPC server listens on the tcp socket connection
	gs.Serve(listener)
}
