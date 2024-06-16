package sqlite_test

import (
	"context"
	"invoicehub"
	"invoicehub/sqlite"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_companyRepository(t *testing.T) {
	db, err := sqlite.SetupDB(":memory:")
	require.NoError(t, err)

	repo := sqlite.NewCompanyRepository(db)
	ctx := context.Background()

	company := invoicehub.Company{
		Name:  "Example Ltd",
		Email: "info@example.com",
		TaxID: "123456789",
		VatID: "987654321",
		Address: invoicehub.Address{
			Street:     "123 Example St",
			City:       "Limassol",
			PostalCode: "12345",
			Country:    "CY",
		},
		BankAccounts: []invoicehub.BankAccount{
			{
				BankName:      "Example Bank",
				AccountNumber: "1234567890",
				IBAN:          "EX12345678901234567890",
				BIC:           "EXBIC123",
			},
		},
	}

	// Create a new company
	id, err := repo.Create(ctx, &company)
	require.NoError(t, err)
	require.NotZero(t, id)
	require.Equal(t, id, company.ID)

	// Get the company
	company2, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, company, company2)

	// Update the company
	company.Name = "Example Ltd 2"
	err = repo.Update(ctx, &company)
	require.NoError(t, err)

	updatedCompany, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, company, updatedCompany)

	// not found
	_, err = repo.Get(ctx, 999)
	require.Error(t, err)
}
