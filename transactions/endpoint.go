package transactions

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints represents all endpoints
type Endpoints struct {
	GetByUserEndpoint endpoint.Endpoint
	GetEndpoint       endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetByUserEndpoint: makeGetByUserEndpoint(s),
		GetEndpoint:       makeGetEndpoint(s),
	}
}

func makeGetByUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransactionsRequest)
		transactions, total, err := s.GetByUser(ctx, *req.UserID, req.Type, req.Sort, req.Page, req.PageSize, req.DateFrom, req.DateTo, req.Open)

		if nil == err {
			return TransactionsResponse{Transactions: transactions, Total: total}, nil
		}

		return TransactionsResponse{}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransactionsRequest)
		transactions, total, err := s.GetByDateRange(ctx, req.Type, req.DateFrom, req.DateTo)

		if nil == err {
			return TransactionsResponse{Transactions: transactions, Total: total}, nil
		}

		return TransactionsResponse{}, err
	}
}
