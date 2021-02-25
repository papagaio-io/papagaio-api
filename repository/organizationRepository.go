package repository

import (
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func (db *AppDb) GetOrganizations() []model.Organization {

	return []model.Organization{
		{OrganizationName: "Test the Database from real database",
			OrganizationType: "Test the Database from real database",
			OrganizationURL:  "Test the Database from real database"},
	}
}
