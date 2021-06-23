package model

import "wecode.sorint.it/opensource/papagaio-api/types"

type GitSource struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	GitType           types.GitType `json:"gitType"`
	GitAPIURL         string        `json:"gitApiUrl"`
	GitClientID       string        `json:"gitClientId"`
	GitSecret         string        `json:"gitSecret"`
	AgolaRemoteSource string        `json:"agolaRemoteSource"`
}
