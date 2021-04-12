package utils

import (
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

type ConfigUtils struct {
	Db repository.Database
}

func (tg *ConfigUtils) GetOrganizationsTriggerTime() uint {
	val := tg.Db.GetOrganizationsTriggerTime()
	if val == -1 {
		return config.Config.TriggersConfig.OrganizationsDefaultTriggerTime
	}
	return uint(val)
}

func (fg *ConfigUtils) GetRunFailedTriggerTime() uint {
	val := fg.Db.GetRunFailedTriggerTime()
	if val == -1 {
		return config.Config.TriggersConfig.RunFailedDefaultTriggerTime
	}
	return uint(val)
}
