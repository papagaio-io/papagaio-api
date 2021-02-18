package database

import "wecode.sorint.it/opensource/papagaio-be/dto"

func (db *Database) GetOrganizations() []dto.Organization {
	//call the real database

	return []dto.Organization{
		{Name: "Test the Database from real database"},
	}
}
