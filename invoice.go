package invoicehub

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *Invoice) (int, error)
	Get(ctx context.Context, id int) (Invoice, error)
	GetLastInvoiceForYear(ctx context.Context, year int) (Invoice, error)
}

type InvoiceService interface {
	Create(ctx context.Context, invoice *Invoice) error
	Get(ctx context.Context, id int) (Invoice, error)
	CreatePDF(ctx context.Context, id int) ([]byte, error)
}

type Invoice struct {
	ID            int
	InvoiceNumber string
	IssueDate     time.Time
	SellerID      int
	BuyerID       int
	DaysToPay     int

	LineItems []LineItem
}

type LineItem struct {
	Description string
	Amount      Amount
	VatRate     decimal.Decimal
}

type Amount struct {
	Value    decimal.Decimal
	Currency string
}
