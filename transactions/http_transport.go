package transactions

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fsilberstein/parameters-issue/errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler ...
func MakeHTTPHandler(endpoints Endpoints, router *mux.Router) http.Handler {

	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorEncoder(errors.LoggingErrorEncoder),
	}

	getByUserHandler := kithttp.NewServer(
		endpoints.GetByUserEndpoint,
		decodeGetByUserRequest,
		encodeResponse,
		options...,
	)

	getHandler := kithttp.NewServer(
		endpoints.GetEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	)

	ur := router.PathPrefix("/users").Subrouter().StrictSlash(true)
	{
		ur.Handle("/{id}/transactions/", getByUserHandler).Methods("GET")
	}

	tr := router.PathPrefix("/transactions").Subrouter().StrictSlash(true)
	{
		tr.Handle("/", getHandler).Methods("GET")
	}

	return router
}

func decodeGetByUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.NewInvalidArgument("invalid parameter 'user_id'")
	}

	//init request with default value and handle non-mandatory parameters
	req := TransactionsRequest{UserID: &id, Sort: "desc"}

	params := r.URL.Query()

	// check `type` parameter
	if transactionType, pTypeOk := params["type"]; pTypeOk && len(transactionType) > 0 {
		var typeArray []string
		for _, value := range transactionType {
			if isTypeValid(value) {
				typeArray = append(typeArray, value)
			} else {
				return nil, errors.NewInvalidArgument("parameter 'type' does not match any of the accepted values")
			}
		}

		if len(typeArray) > 0 {
			req.Type = typeArray
		}
	}

	sort, ok := params["sort"]
	if ok && len(sort) > 0 {
		// right now, only asc and desc are accepted
		if sort[0] != "asc" && sort[0] != "desc" {
			return nil, errors.NewInvalidArgument("invalid parameter 'sort'")
		}
		req.Sort = sort[0]
	}

	page, ok := params["page"]
	if ok && len(page) > 0 {
		i, err := strconv.ParseInt(page[0], 10, 0)
		if err != nil || i < 1 {
			return nil, errors.NewInvalidArgument("invalid parameter 'page'")
		}
		req.Page = int(i)
	} else {
		req.Page = 1
	}

	pageSize, ok := params["page_size"]
	if ok && len(pageSize) > 0 {
		i, err := strconv.ParseInt(pageSize[0], 10, 0)
		if err != nil || i < 1 {
			return nil, errors.NewInvalidArgument("invalid parameter 'page_size'")
		}
		req.PageSize = int(i)
	} else {
		req.PageSize = -1
	}

	dateFromStr, ok := params["date_from"]
	if ok {
		dateFrom, err := time.Parse(time.RFC3339, dateFromStr[0])
		if err != nil {
			return nil, errors.NewInvalidArgument("could not decode `date_from`")
		}
		req.DateFrom = &dateFrom
	}

	dateToStr, ok := params["date_to"]
	if ok {
		dateTo, err := time.Parse(time.RFC3339, dateToStr[0])
		if err != nil {
			return nil, errors.NewInvalidArgument("could not decode `date_to`")
		}
		req.DateTo = &dateTo
	}

	open, ok := params["open"]
	if ok && len(open) > 0 {
		o, err := strconv.ParseBool(open[0])
		if err != nil {
			return nil, errors.NewInvalidArgument("invalid parameter 'open'")
		}
		req.Open = &o
	}

	return req, nil
}

func decodeGetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params := r.URL.Query()

	request := TransactionsRequest{}

	// check `type` parameter
	if transactionType, pTypeOk := params["type"]; pTypeOk && len(transactionType) > 0 {
		var typeArray []string
		for _, value := range transactionType {
			if isTypeValid(value) {
				typeArray = append(typeArray, value)
			} else {
				return nil, errors.NewInvalidArgument("parameter 'type' does not match any of the accepted values")
			}
		}

		if len(typeArray) > 0 {
			request.Type = typeArray
		}
	}

	dateFromStr, ok := params["date_from"]
	if ok {
		dateFrom, err := time.Parse(time.RFC3339, dateFromStr[0])
		if err != nil {
			return nil, errors.NewInvalidArgument("could not decode `date_from`")
		}
		request.DateFrom = &dateFrom
	}

	dateToStr, ok := params["date_to"]
	if ok {
		dateTo, err := time.Parse(time.RFC3339, dateToStr[0])
		if err != nil {
			return nil, errors.NewInvalidArgument("could not decode `date_to`")
		}
		request.DateTo = &dateTo
	}

	if request.DateFrom == nil && request.DateTo == nil {
		return nil, errors.NewInvalidArgument("at least one of the date range boundaries must be set")
	}

	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func isTypeValid(transactionType string) bool {
	return true
}
