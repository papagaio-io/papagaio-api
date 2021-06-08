package dto

import (
	"errors"
	"net/url"

	"wecode.sorint.it/opensource/papagaio-api/types"
)

type GitSourcesDto struct {
	Name      string `json:"name"`
	GitAPIURL string `json:"gitApiUrl"`
}

type UpdateGitSourceRequestDto struct {
	GitType           *types.GitType `json:"gitType"`
	GitAPIURL         *string        `json:"gitApiUrl"`
	AgolaRemoteSource *string        `json:"agolaRemoteSource"`
}

type CreateGitSourceDto struct {
	Name      string        `json:"name"`
	GitType   types.GitType `json:"gitType"`
	GitAPIURL string        `json:"gitApiUrl"`

	GitClientID string `json:"gitClientId"`
	GitSecret   string `json:"gitSecret"`

	AgolaRemoteSource *string `json:"agolaRemoteSource"`
	AgolaClientID     *string `json:"agolaClientId"`
	AgolaSecret       *string `json:"agolaSecret"`
}

func (gitSource *CreateGitSourceDto) IsValid() error {
	if len(gitSource.Name) == 0 {
		return errors.New("name is empty")
	}

	err := gitSource.GitType.IsValid()
	if err != nil {
		return err

	}

	_, err = url.ParseRequestURI(gitSource.GitAPIURL)
	if err != nil {
		return errors.New("gitApiUrl is not valid")
	}

	if len(gitSource.GitClientID) == 0 {
		return errors.New("gitClientId is empty")
	}

	if len(gitSource.GitSecret) == 0 {
		return errors.New("gitSecret is empty")
	}

	if gitSource.AgolaRemoteSource == nil {
		if gitSource.AgolaClientID == nil || gitSource.AgolaSecret == nil {
			return errors.New("agolaRemoteSource or oauth2 application must be specified")
		}
	}

	return nil
}
