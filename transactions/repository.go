package transactions

import (
	"context"
	"time"
)

// Repository interface
type Repository interface {
	GetByUser(ctx context.Context, userID string, transactionType []string, sort string, page, pageSize int, dateFrom, dateTo *time.Time, open *bool) ([]*Transaction, int64, error)
	GetByDateRange(ctx context.Context, transactionType []string, dateFrom, dateTo *time.Time) ([]*Transaction, int64, error)
}
