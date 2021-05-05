package dto

import "wecode.sorint.it/opensource/papagaio-api/types"

type GitSourcesDto struct {
	Name      string `json:"name"`
	GitAPIURL string `json:"gitApiUrl"`
}

type UpdateRemoteSourceRequestDto struct {
	GitType           *types.GitType `json:"gitType"`
	GitAPIURL         *string        `json:"gitApiUrl"`
	GitToken          *string        `json:"gitToken"`
	AgolaRemoteSource *string        `json:"agolaRemoteSource"`
	AgolaToken        *string        `json:"agolaToken"`
}
