package service

import (
	"wecode.sorint.it/opensource/papagaio-be/database"
	"wecode.sorint.it/opensource/papagaio-be/dto"
)

type ServiceInterface interface {
	GetOrganizations() []dto.Organization
}

type Service struct {
	Db database.DatabaseInterface
}
