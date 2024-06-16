package sqlite

import (
	"context"
	"invoicehub"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var _ invoicehub.CompanyRepository = (*companyRepository)(nil)

type dbcompany struct {
	ID           int    `gorm:"primaryKey"`
	Name         string `gorm:"type:text"`
	Address      datatypes.JSONType[invoicehub.Address]
	Email        string `gorm:"type:text"`
	TaxID        string `gorm:"type:text"`
	VatID        string `gorm:"type:text"`
	BankAccounts datatypes.JSONSlice[invoicehub.BankAccount]
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) invoicehub.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(ctx context.Context, company *invoicehub.Company) (int, error) {
	// Convert invoicehub.Company to dbcompany
	dbCompany := dbcompany{
		Name:         company.Name,
		Address:      datatypes.NewJSONType(company.Address),
		Email:        company.Email,
		TaxID:        company.TaxID,
		VatID:        company.VatID,
		BankAccounts: datatypes.NewJSONSlice(company.BankAccounts),
	}

	// Save to database

	if err := r.db.WithContext(ctx).Create(&dbCompany).Error; err != nil {
		return 0, err
	}

	company.ID = dbCompany.ID

	return company.ID, nil
}

func (r *companyRepository) Update(ctx context.Context, company *invoicehub.Company) error {
	var existingCompany dbcompany
	if err := r.db.WithContext(ctx).First(&existingCompany, company.ID).Error; err != nil {
		return err
	}

	existingCompany.Name = company.Name
	existingCompany.Address = datatypes.NewJSONType(company.Address)
	existingCompany.Email = company.Email
	existingCompany.TaxID = company.TaxID
	existingCompany.VatID = company.VatID
	existingCompany.BankAccounts = datatypes.NewJSONSlice(company.BankAccounts)

	if err := r.db.WithContext(ctx).Save(&existingCompany).Error; err != nil {
		return err
	}

	return nil
}

func (r *companyRepository) Get(ctx context.Context, id int) (invoicehub.Company, error) {
	var dbCompany dbcompany
	if err := r.db.WithContext(ctx).First(&dbCompany, id).Error; err != nil {
		return invoicehub.Company{}, err
	}

	company := invoicehub.Company{
		ID:           dbCompany.ID,
		Name:         dbCompany.Name,
		Address:      dbCompany.Address.Data(),
		Email:        dbCompany.Email,
		TaxID:        dbCompany.TaxID,
		VatID:        dbCompany.VatID,
		BankAccounts: dbCompany.BankAccounts,
	}

	return company, nil
}
