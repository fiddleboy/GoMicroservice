package server

import (
	"context"
	"gRPC-play/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	log hclog.Logger
	currency.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger) Currency {
	return Currency{l, currency.UnimplementedCurrencyServer{}}
}

func (c Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GeRate", "base", rr.GetBase(), rr.GetDestination())

	return &currency.RateResponse{Rate: 0.5}, nil
}
