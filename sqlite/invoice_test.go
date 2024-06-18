package sqlite_test

import (
	"context"
	"invoicehub"
	"invoicehub/sqlite"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func compareInvoice(t *testing.T, a, b invoicehub.Invoice) {
	require.Equal(t, a.ID, b.ID, "ID mismatch")
	require.Equal(t, a.InvoiceNumber, b.InvoiceNumber, "InvoiceNumber mismatch")
	require.Equal(t, a.IssueDate, b.IssueDate, "IssueDate mismatch")
	require.Equal(t, a.SellerID, b.SellerID, "SellerID mismatch")
	require.Equal(t, a.BuyerID, b.BuyerID, "BuyerID mismatch")
	require.Equal(t, a.DaysToPay, b.DaysToPay, "DaysToPay mismatch")

	require.Len(t, a.LineItems, len(b.LineItems), "LineItems length mismatch")

	for i := range a.LineItems {
		require.Equal(t, a.LineItems[i].Description, b.LineItems[i].Description, "Description mismatch")
		require.Equal(t, a.LineItems[i].Amount.Currency, b.LineItems[i].Amount.Currency, "Currency mismatch")
		require.Equal(t, a.LineItems[i].Amount.Value.String(), b.LineItems[i].Amount.Value.String(), "Amount value mismatch")
	}
}

func Test_invoiceRepository(t *testing.T) {
	db, err := sqlite.SetupDB(":memory:")
	require.NoError(t, err)

	repo := sqlite.NewInvoiceRepository(db)
	require.NotNil(t, repo)

	ctx := context.Background()

	invoice := invoicehub.Invoice{
		InvoiceNumber: "1001/24",
		IssueDate:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		SellerID:      1,
		BuyerID:       2,
		DaysToPay:     14,

		LineItems: []invoicehub.LineItem{
			{
				Description: "My awesome services",
				Amount: invoicehub.Amount{
					Value:    decimal.NewFromFloat(100.00),
					Currency: "EUR",
				},
				VatRate: decimal.NewFromFloat(0.19),
			},
		},
	}

	// Create a new invoice
	id, err := repo.Create(ctx, &invoice)
	require.NoError(t, err)
	require.NotZero(t, id)
	require.Equal(t, id, invoice.ID)

	// fetch the invoice

	invoice2, err := repo.Get(ctx, id)
	require.NoError(t, err)
	compareInvoice(t, invoice, invoice2)

	// fetch a non-existing invoice
	_, err = repo.Get(ctx, 999)
	require.Error(t, err)

	// add one more invoice
	invoice3 := invoicehub.Invoice{
		InvoiceNumber: "1002/24",
		IssueDate:     time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
		SellerID:      1,
		BuyerID:       2,
		DaysToPay:     14,

		LineItems: []invoicehub.LineItem{
			{
				Description: "web development",
				Amount: invoicehub.Amount{
					Value:    decimal.NewFromFloat(100.00),
					Currency: "EUR",
				},
				VatRate: decimal.NewFromFloat(0.19),
			},
			{
				Description: "web scraping",
				Amount: invoicehub.Amount{
					Value:    decimal.NewFromFloat(150.00),
					Currency: "EUR",
				},
				VatRate: decimal.NewFromFloat(0.19),
			},
		},
	}

	id3, err := repo.Create(ctx, &invoice3)
	require.NoError(t, err)
	require.NotZero(t, id3)
	require.Equal(t, id3, invoice3.ID)

	// fetch the last invoice for the year 2024

	lastInvoice, err := repo.GetLastInvoiceForYear(ctx, 2024)
	require.NoError(t, err)
	compareInvoice(t, invoice3, lastInvoice)
}
