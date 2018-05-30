package transactions

import (
	"context"
	"time"
)

// Service is the transaction service interface
type Service interface {
	GetByUser(ctx context.Context, userID string, transactionType []string, sort string, page, pageSize int, dateFrom, dateTo *time.Time, open *bool) ([]*Transaction, int64, error)
	GetByDateRange(ctx context.Context, transactionType []string, dateFrom, dateTo *time.Time) ([]*Transaction, int64, error)
}

type service struct {
	repo Repository
}

// NewService initializes new service
func NewService(repo Repository) (Service, error) {
	return &service{
		repo: repo,
	}, nil
}

func (s *service) GetByUser(ctx context.Context, userID string, transactionType []string, sort string, page, pageSize int, dateFrom, dateTo *time.Time, open *bool) ([]*Transaction, int64, error) {
	return s.repo.GetByUser(ctx, userID, transactionType, sort, page, pageSize, dateFrom, dateTo, open)
}

func (s *service) GetByDateRange(ctx context.Context, transactionType []string, dateFrom, dateTo *time.Time) ([]*Transaction, int64, error) {
	return s.repo.GetByDateRange(ctx, transactionType, dateFrom, dateTo)
}
