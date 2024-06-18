package sqlite

import (
	"context"
	"invoicehub"
	"strconv"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var _ invoicehub.InvoiceRepository = &invoiceRepository{}

type dbinvoice struct {
	ID            int       `gorm:"primaryKey"`
	InvoiceNumber string    `gorm:"type:text"`
	IssueDate     time.Time `gorm:"type:datetime"`
	BuyerID       int
	SellerID      int
	DaysToPay     int

	LineItems datatypes.JSONSlice[invoicehub.LineItem]
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) invoicehub.InvoiceRepository {
	return &invoiceRepository{db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *invoicehub.Invoice) (int, error) {
	dbInvoice := dbinvoice{
		InvoiceNumber: invoice.InvoiceNumber,
		IssueDate:     invoice.IssueDate,
		BuyerID:       invoice.BuyerID,
		SellerID:      invoice.SellerID,
		DaysToPay:     invoice.DaysToPay,
		LineItems:     datatypes.NewJSONSlice(invoice.LineItems),
	}

	if err := r.db.WithContext(ctx).Create(&dbInvoice).Error; err != nil {
		return 0, err
	}

	invoice.ID = dbInvoice.ID

	return invoice.ID, nil
}

func (r *invoiceRepository) Get(ctx context.Context, id int) (invoicehub.Invoice, error) {
	var dbitem dbinvoice
	if err := r.db.WithContext(ctx).First(&dbitem, id).Error; err != nil {
		return invoicehub.Invoice{}, err
	}

	invoice := invoicehub.Invoice{
		ID:            dbitem.ID,
		InvoiceNumber: dbitem.InvoiceNumber,
		IssueDate:     dbitem.IssueDate,
		BuyerID:       dbitem.BuyerID,
		SellerID:      dbitem.SellerID,
		DaysToPay:     dbitem.DaysToPay,
		LineItems:     dbitem.LineItems,
	}

	return invoice, nil
}

func (r *invoiceRepository) GetLastInvoiceForYear(ctx context.Context, year int) (invoicehub.Invoice, error) {
	query := r.db.WithContext(ctx).
		Where("strftime('%Y', issue_date) = ?", strconv.Itoa(year)).
		Order("issue_date desc").
		Limit(1)

	var dbitem dbinvoice
	if err := query.First(&dbitem).Error; err != nil {
		return invoicehub.Invoice{}, err
	}

	invoice := invoicehub.Invoice{
		ID:            dbitem.ID,
		InvoiceNumber: dbitem.InvoiceNumber,
		IssueDate:     dbitem.IssueDate,
		BuyerID:       dbitem.BuyerID,
		SellerID:      dbitem.SellerID,
		DaysToPay:     dbitem.DaysToPay,
		LineItems:     dbitem.LineItems,
	}

	return invoice, nil
}
