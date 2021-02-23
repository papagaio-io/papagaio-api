package database

import (
	"wecode.sorint.it/opensource/papagaio-be/dto"
)

func (db *Database) GetOrganizations() []dto.Organization {

	return []dto.Organization{
		{OrganizationName: "Test the Database from real database",
			OrganizationType: "Test the Database from real database",
			OrganizationURL:  "Test the Database from real database"},
	}
}
