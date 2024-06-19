package invoices

import (
	"context"
	"errors"
	"fmt"
	"invoicehub"
	"strconv"
	"strings"
)

var _ invoicehub.InvoiceService = (*invoiceSvc)(nil)

type invoiceSvc struct {
	repo invoicehub.InvoiceRepository
}

func New(repo invoicehub.InvoiceRepository) invoicehub.InvoiceService {
	return &invoiceSvc{repo}
}

func (svc *invoiceSvc) Create(ctx context.Context, invoice *invoicehub.Invoice) error {
	issueYear := invoice.IssueDate.Year()

	var invoiceNumber string

	lastInvoice, err := svc.repo.GetLastInvoiceForYear(ctx, issueYear)
	if err != nil {
		if errors.Is(err, invoicehub.ErrInvoiceNotFound) {
			invoiceNumber = fmt.Sprintf("1001/%d", issueYear)
		} else {
			return err
		}
	}

	if invoiceNumber == "" {
		parts := strings.Split(lastInvoice.InvoiceNumber, "/")
		if len(parts) != 2 {
			return errors.New("invalid invoice number")
		}

		lastInvoiceNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		invoiceNumber = fmt.Sprintf("%d/%d", lastInvoiceNumber+1, issueYear)
	}

	invoice.InvoiceNumber = invoiceNumber

	if _, err := svc.repo.Create(ctx, invoice); err != nil {
		return err
	}

	return nil
}

func (svc *invoiceSvc) Get(ctx context.Context, id int) (invoicehub.Invoice, error) {
	return svc.repo.Get(ctx, id)
}

func (svc *invoiceSvc) CreatePDF(ctx context.Context, id int) ([]byte, error) {
	return nil, nil
}
