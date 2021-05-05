package model

import "wecode.sorint.it/opensource/papagaio-api/types"

type GitSource struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	GitType           types.GitType `json:"gitType"`
	GitAPIURL         string        `json:"gitApiUrl"`
	GitToken          string        `json:"gitToken"`
	AgolaRemoteSource string        `json:"agolaRemoteSource"`
	AgolaToken        string        `json:"agolaToken"`
}
