package invoices_test

import (
	"context"
	"errors"
	"invoicehub"
	"invoicehub/invoices"
	"invoicehub/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_New(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	repo := mocks.NewMockInvoiceRepository(mctrl)

	svc := invoices.New(repo)
	require.NotNil(t, svc)
}

func Test_invoiceSvc_Create(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	repo := mocks.NewMockInvoiceRepository(mctrl)

	svc := invoices.New(repo)

	t.Run("when there is no last invoice for the current year", func(t *testing.T) {
		inv := &invoicehub.Invoice{
			IssueDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		repo.EXPECT().GetLastInvoiceForYear(gomock.Any(), 2024).Return(invoicehub.Invoice{}, invoicehub.ErrInvoiceNotFound)
		repo.EXPECT().Create(gomock.Any(), inv).Return(1, nil)

		err := svc.Create(context.Background(), inv)
		require.NoError(t, err)

		require.Equal(t, "1001/2024", inv.InvoiceNumber)
	})

	t.Run("when there is a last invoice for the current year", func(t *testing.T) {
		inv1 := &invoicehub.Invoice{
			IssueDate:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			InvoiceNumber: "1001/2024",
		}

		invNew := &invoicehub.Invoice{
			IssueDate: time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
		}

		repo.EXPECT().GetLastInvoiceForYear(gomock.Any(), 2024).Return(*inv1, nil)

		repo.EXPECT().Create(gomock.Any(), invNew).Return(2, nil)

		err := svc.Create(context.Background(), invNew)
		require.NoError(t, err)

		require.Equal(t, "1002/2024", invNew.InvoiceNumber)
	})

	t.Run("when there is a database error while getting the last invoice", func(t *testing.T) {
		inv := &invoicehub.Invoice{
			IssueDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		repo.EXPECT().GetLastInvoiceForYear(gomock.Any(), 2024).Return(invoicehub.Invoice{}, errors.New("something went wrong"))

		err := svc.Create(context.Background(), inv)
		require.Error(t, err)
	})

	t.Run("when we cannot save the invoice", func(t *testing.T) {
		invoice := &invoicehub.Invoice{
			IssueDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		repo.EXPECT().GetLastInvoiceForYear(gomock.Any(), 2024).Return(invoicehub.Invoice{}, invoicehub.ErrInvoiceNotFound)
		repo.EXPECT().Create(gomock.Any(), invoice).Return(0, errors.New("something went wrong"))

		err := svc.Create(context.Background(), invoice)
		require.Error(t, err)
	})
}

func Test_invoiceSvc_Get(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	repo := mocks.NewMockInvoiceRepository(mctrl)

	svc := invoices.New(repo)

	t.Run("when the invoice is found", func(t *testing.T) {
		expected := invoicehub.Invoice{
			ID:            1,
			InvoiceNumber: "1001/2024",
			IssueDate:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			SellerID:      1,
			BuyerID:       2,
			DaysToPay:     14,
		}

		repo.EXPECT().Get(gomock.Any(), 1).Return(expected, nil)

		inv, err := svc.Get(context.Background(), 1)
		require.NoError(t, err)

		require.Equal(t, expected, inv)
	})

	t.Run("when the invoice is not found", func(t *testing.T) {
		repo.EXPECT().Get(gomock.Any(), 1).Return(invoicehub.Invoice{}, invoicehub.ErrInvoiceNotFound)

		_, err := svc.Get(context.Background(), 1)
		require.Error(t, err)
		require.ErrorIs(t, err, invoicehub.ErrInvoiceNotFound)
	})

	t.Run("when there is a database error", func(t *testing.T) {
		repo.EXPECT().Get(gomock.Any(), 1).Return(invoicehub.Invoice{}, errors.New("something went wrong"))

		_, err := svc.Get(context.Background(), 1)
		require.Error(t, err)
	})
}

func Test_invoiceSvc_CreatePDF(t *testing.T) {
}
