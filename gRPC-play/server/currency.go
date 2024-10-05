package server

import (
	"context"
	"io"
	"time"

	"github.com/fiddleboy/GoMicroservice/gRPC-play/data"
	"github.com/fiddleboy/GoMicroservice/gRPC-play/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	log           hclog.Logger
	rates         *data.ExchangeRates
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
	currency.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger, r *data.ExchangeRates) *Currency {
	c := Currency{l, r, make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest), currency.UnimplementedCurrencyServer{}}
	// start a thread to keep monitoring rate changes and notify subscribed clients
	go c.handleUpdates()
	return &c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		// c.log.Info("Got updated rates!")

		for k, v := range c.subscriptions {
			for _, rateReq := range v {
				rate, err := c.rates.GetRate(rateReq.Base.String(), rateReq.Destination.String())
				if err != nil {
					c.log.Error("Unable to get updated rate", "base", rateReq.GetBase().String(), "destination", rateReq.Destination.String())
					continue
				}
				// send the new rate to the subscribed client
				err = k.Send(&currency.RateResponse{Rate: rate})
				if err != nil {
					// TOOD: need to check if client has already losed the conn,
					// if so, then we need remove this subscription from our table
					delete(c.subscriptions, k)
					c.log.Error("Unable to send updated rate", "base", rateReq.GetBase().String(), "destination", rateReq.Destination.String())
				}

			}
		}
	}
}

func (c Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.Base, "dest", rr.Destination)
	if rr.Base == rr.Destination {
		errStatus := status.Newf(
			codes.InvalidArgument,
			"Base and Dest are the same!",
		)
		enrichedErrStatus, err_ := errStatus.WithDetails(rr)
		if err_ != nil {
			return nil, err_
		}
		return nil, enrichedErrStatus.Err()
	}
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	c.log.Info("The rate retrieved is: ", rate)
	return &currency.RateResponse{Rate: rate}, nil
}

func (c Currency) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {
	// for {
	// 	err := src.Send(&currency.RateResponse{Rate: 12.1})
	// 	if err != nil {
	// 		// c.log.Error("After receiving end screwed, sending end also exited after the err check!\n", "The err checked on the sending side is: ", err)
	// 		return err
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }

	for {
		rr, err := src.Recv()
		if err == io.EOF {
			c.log.Info("Client has closed connection.")
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			break
		}
		c.log.Info("Handle client request", "request", rr)
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*currency.RateRequest{}
		}

		var validationErr *status.Status

		for _, v := range rrs {
			if v.Base == rr.Base || v.Destination == rr.Destination {
				validationErr = status.Newf(
					codes.AlreadyExists,
					"Unable to subscirbe for already existed currency.",
				)
			}
			validationErr, err = validationErr.WithDetails(rr)
			if err != nil {
				c.log.Error("Unable to add metadata to the status", "error", err)
				break
			}
		}

		if validationErr != nil {
			src.Send(
				&currency.StreamingRateResponse{
					Message: &currency.StreamingRateResponse_Error{
						Error: validationErr.Proto(),
					},
				},
			)
			continue
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
		c.log.Info("The curren subscription table is: ", src, rrs)
	}
	return nil
}
