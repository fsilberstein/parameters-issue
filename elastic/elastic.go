package elastic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/facebookgo/httpcontrol"
	"github.com/fsilberstein/parameters-issue/logger"
	loghttp "github.com/motemen/go-loghttp"
	"github.com/motemen/go-nuts/roundtime"
	elasticapi "gopkg.in/olivere/elastic.v5"
)

var (
	// ErrElasticSearchNotReachable ...
	ErrElasticSearchNotReachable = errors.New("Elastic search server is not reachable")
	elasticResponseSize          = -1
)

const (
	DocumentTypeFee         = "fee"
	DocumentTypeCredit      = "credit"
	DocumentTypeRefund      = "refund"
	DocumentTypePayment     = "payment"
	DocumentTypeUser        = "user"
	DocumentTypeInvoice     = "invoice"
	DocumentTypeTransaction = "transaction"
	DocumentTypeReceipt     = "receipt"
	DocumentTypeBalance     = "balance"

	timeout = 2 * time.Second
)

// createHTTPClient ...
func createHTTPClient(isElasticDebug bool) *http.Client {
	// by default use retry system with timeout
	httpClient := &http.Client{Transport: &httpcontrol.Transport{
		MaxTries:          3,
		RequestTimeout:    timeout,
		RetryAfterTimeout: true,
	}}

	// if we are in debug mode, we can skip the retry and only print request and response
	if isElasticDebug {
		httpClient = &http.Client{
			Timeout: timeout,
			Transport: &loghttp.Transport{
				LogRequest: func(req *http.Request) {

					var bodyBuffer []byte
					if req.Body != nil {
						bodyBuffer, _ = ioutil.ReadAll(req.Body) // after this operation body will equal 0
						// Restore the io.ReadCloser to request
						req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
					}

					fmt.Println("--------- Elasticsearch ---------")
					fmt.Println("Request URL : ", req.URL)
					fmt.Println("Request Method : ", req.Method)
					fmt.Println("Request Body : ", string(bodyBuffer))

				},
				LogResponse: func(resp *http.Response) {
					ctx := resp.Request.Context()
					if start, ok := ctx.Value(loghttp.ContextKeyRequestStart).(time.Time); ok {
						fmt.Println("Response Status : ", resp.StatusCode)
						fmt.Println("Response Duration : ", roundtime.Duration(time.Now().Sub(start), 2))
					} else {
						fmt.Println("Response Status : ", resp.StatusCode)
					}
					fmt.Println("--------------------------------")
				},
			},
		}
	}

	return httpClient
}

// NewElasticClient ...
func NewElasticClient(ctx context.Context, url string, sniff bool, responseSize int, isElasticDebug bool) (*elasticapi.Client, error) {
	elasticResponseSize = responseSize

	httpClient := createHTTPClient(isElasticDebug)
	client, err := elasticapi.NewClient(elasticapi.SetURL(url), elasticapi.SetSniff(sniff), elasticapi.SetHttpClient(httpClient))

	if err != nil {
		return nil, err
	}

	err = ping(ctx, client, url)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Ping method
func ping(ctx context.Context, client *elasticapi.Client, url string) error {

	// Ping the Elasticsearch server to get HttpStatus, version number
	if client != nil {
		info, code, err := client.Ping(url).Do(ctx)
		if err != nil {
			return err
		}

		logger.LogStdOut.Info(fmt.Sprintf("Elasticsearch returned with code %d and version %s", code, info.Version.Number))
		return nil
	}

	return errors.New("elastic client is nil")
}

func getSort(sort string) *elasticapi.FieldSort {
	if sort == "asc" {
		return elasticapi.NewFieldSort("creation_date").Asc()
	}
	return elasticapi.NewFieldSort("creation_date").Desc()
}
