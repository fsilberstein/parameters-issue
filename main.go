package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fsilberstein/parameters-issue/config"
	"github.com/fsilberstein/parameters-issue/elastic"
	"github.com/fsilberstein/parameters-issue/logger"
	"github.com/fsilberstein/parameters-issue/transactions"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	appName = "bookkeeping-transaction-viewer"
)

var (
	err     error
	errc    chan error
	brokers []string
)

func init() {
	errc = make(chan error, 10)
}

func main() {
	ctx := context.Background()

	// Interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// elasticClient init
	elasticClient, err := elastic.NewElasticClient(ctx, config.ElasticHost, config.ElasticSniff, config.ElasticResponseSize, config.ElasticDebug)
	if err != nil {
		logger.LogStdErr.Error(err)
	}

	// Creates transactions service
	var transactionsService transactions.Service
	var transactionRepository transactions.Repository
	{
		transactionRepository = elastic.NewTransactionRepository(config.ElasticIndex, elasticClient)
		transactionsService, err = transactions.NewService(transactionRepository)
		if err != nil {
			logger.LogStdErr.Error(err)
		}
	}

	// Transaction endpoint
	transactionsEndpoint := transactions.MakeEndpoints(transactionsService)

	// Instances a new HTTP server for healthy check and metrics
	go func() {
		httpAddr := ":" + strconv.Itoa(config.Port)
		mux := mux.NewRouter()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Welcome to the my problem API!")
		})
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Init and register to the router the various endpoints
		transactions.MakeHTTPHandler(transactionsEndpoint, mux)

		logger.LogStdOut.Info(fmt.Sprintf("The API is started on port %d", config.Port))
		errc <- http.ListenAndServe(httpAddr, mux)
	}()

	logger.LogStdErr.Error(<-errc)
}
