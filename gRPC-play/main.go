package main

import (
	"gRPC-play/protos/currency"
	"gRPC-play/server"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)
	currency.RegisterCurrencyServer(gs, cs)
	reflection.Register(gs)

	// set up a tcp connection
	l, err := net.Listen("tcp", ":9092")

	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	// have the gRPC server listens on the tcp socket connection
	gs.Serve(l)
}
