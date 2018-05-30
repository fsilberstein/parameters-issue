package transactions

import (
	"time"
)

type TransactionsRequest struct {
	UserID   *string    `json:"user_id"`
	Sort     string     `json:"sort"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Type     []string   `json:"type"`
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
	Open     *bool      `json:"open"`
}

type TransactionsResponse struct {
	Transactions []*Transaction `json:"transactions"`
	Total        int64          `json:"total"`
}

// Transaction struct (Just use Only one field)
type Transaction struct {
	ID string `json:"id,omitempty"`
}
