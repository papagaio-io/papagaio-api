package service

import "wecode.sorint.it/opensource/papagaio-be/dto"

func (service *Service) GetOrganizations() []dto.Organization {
	return service.Db.GetOrganizations()

}
