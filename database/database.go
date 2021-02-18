package database

import "wecode.sorint.it/opensource/papagaio-be/dto"

type DatabaseInterface interface {
	GetOrganizations() []dto.Organization
}

type Database struct {
}
