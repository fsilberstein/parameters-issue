package elastic

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/fsilberstein/parameters-issue/config"
	apierror "github.com/fsilberstein/parameters-issue/errors"
	"github.com/fsilberstein/parameters-issue/transactions"
	"github.com/pkg/errors"
	elasticapi "gopkg.in/olivere/elastic.v5"
)

type transactionRepository struct {
	IndexName     string
	elasticClient *elasticapi.Client
}

// NewTransactionRepository ...
func NewTransactionRepository(indexName string, elasticClient *elasticapi.Client) transactions.Repository {
	return &transactionRepository{
		IndexName:     indexName,
		elasticClient: elasticClient,
	}
}

// getFromAndSize translates the fields provided by the caller ("page" and "page_size") into
// query parameters that ElasticSearch can understand ("from" and "size")
func getFromAndSize(pageSize int, page int) (int, int, error) {

	var from, size int

	if pageSize < 1 {
		size = elasticResponseSize
	} else if pageSize > elasticResponseSize {
		return 0, 0, apierror.NewInvalidArgument(fmt.Sprintf("invalid parameter 'page_size' > %v", config.ElasticResponseSize))
	} else {
		size = pageSize
	}

	from = (page - 1) * size

	return from, size, nil
}

// GetTransactions ...
func (repo *transactionRepository) GetByUser(ctx context.Context, userID string, transactionType []string, sort string, page, pageSize int, dateFrom, dateTo *time.Time, open *bool) (result []*transactions.Transaction, total int64, err error) {
	if repo.elasticClient == nil {
		err = ErrElasticSearchNotReachable
		return
	}

	// Scope query to a user
	musts := []elasticapi.Query{elasticapi.NewTermQuery("user_id", userID)}

	typeQuery := getTypeQuery(transactionType)
	if typeQuery != nil {
		musts = append(musts, typeQuery)
	}

	dateRangeQuery := getRangeQuery(dateFrom, dateTo)
	if dateRangeQuery != nil {
		musts = append(musts, dateRangeQuery)
	}

	query := elasticapi.NewBoolQuery().Must(musts...)

	from, size, err := getFromAndSize(pageSize, page)
	if err != nil {
		return
	}

	sortObj := getSort(sort)
	searchService := repo.elasticClient.Search(repo.IndexName).
		Index(repo.IndexName). // search in index
		Type(DocumentTypeTransaction).
		Query(query). // specify the query
		SortBy(sortObj).
		Size(size).
		From(from)

	// do request to ElasticSearch
	var searchResult *elasticapi.SearchResult
	searchResult, err = searchService.Do(ctx)
	if err != nil {
		return result, total, errors.Wrap(err, "error during elastic search")
	}

	// then here is the where business logic go to map what is in ES and what we return
	// not relevant here
	total = searchResult.TotalHits()

	return result, total, nil
}

func (repo *transactionRepository) GetByDateRange(ctx context.Context, transactionType []string, dateFrom, dateTo *time.Time) (result []*transactions.Transaction, total int64, err error) {
	if repo.elasticClient == nil {
		err = ErrElasticSearchNotReachable
		return
	}

	var musts []elasticapi.Query

	typeQuery := getTypeQuery(transactionType)
	if typeQuery != nil {
		musts = append(musts, typeQuery)
	}

	dateRangeQuery := getRangeQuery(dateFrom, dateTo)
	if dateRangeQuery != nil {
		musts = append(musts, dateRangeQuery)
	}

	query := elasticapi.NewBoolQuery().Must(musts...)
	scroll := repo.elasticClient.
		Scroll(repo.IndexName).
		Type(DocumentTypeTransaction).
		Query(query).
		Size(elasticResponseSize)

	for {
		searchResult, err := scroll.Do(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, int64(0), err
		}

		// then here is the where business logic go to map what is in ES and what we return
		// not relevant here
		total = searchResult.TotalHits()
	}
	return result, total, nil
}

func getTypeQuery(types []string) *elasticapi.BoolQuery {
	if types != nil {
		var existQueries []elasticapi.Query
		for _, value := range types {
			existQueries = append(existQueries, elasticapi.NewExistsQuery(value))
		}

		return elasticapi.NewBoolQuery().Should(existQueries...)
	}
	return nil
}

func getRangeQuery(dateFrom, dateTo *time.Time) *elasticapi.RangeQuery {
	if dateFrom != nil || dateFrom != nil {
		dateRangeQuery := elasticapi.NewRangeQuery("creation_date").IncludeUpper(false).IncludeLower(true)
		if dateFrom != nil {
			dateRangeQuery.From(*dateFrom)
		}
		if dateTo != nil {
			dateRangeQuery.To(*dateTo)
		}
		return dateRangeQuery
	}
	return nil
}
