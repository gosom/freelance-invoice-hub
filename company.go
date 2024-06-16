package invoicehub

import "context"

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) (int, error)
	Update(ctx context.Context, company *Company) error
	Get(ctx context.Context, id int) (Company, error)
}

type Company struct {
	ID           int
	Name         string
	Address      Address
	Email        string
	TaxID        string
	VatID        string
	BankAccounts []BankAccount
}

type Address struct {
	Street     string
	City       string
	PostalCode string
	Country    string
}

type BankAccount struct {
	BankName      string
	AccountNumber string
	IBAN          string
	BIC           string
}
