package dto

import (
	"errors"
	"net/url"

	"wecode.sorint.it/opensource/papagaio-api/types"
)

type GitSourcesDto struct {
	Name      string        `json:"name"`
	GitAPIURL string        `json:"gitApiUrl"`
	LoginURL  string        `json:"loginUrl"`
	GitType   types.GitType `json:"gitType"`
}

type UpdateGitSourceRequestDto struct {
	GitType         *types.GitType `json:"gitType"`
	GitAPIURL       *string        `json:"gitApiUrl"`
	GitClientID     *string        `json:"gitClientId"`
	GitClientSecret *string        `json:"gitClientSecret"`

	AgolaRemoteSource *string `json:"agolaRemoteSource"`
}

type CreateGitSourceRequestDto struct {
	Name      string        `json:"name"`
	GitType   types.GitType `json:"gitType"`
	GitAPIURL *string       `json:"gitApiUrl"`

	GitClientID     string `json:"gitClientId"`
	GitClientSecret string `json:"gitClientSecret"`

	AgolaRemoteSourceName *string `json:"agolaRemoteSourceName"`
	AgolaClientID         *string `json:"agolaClientId"`
	AgolaClientSecret     *string `json:"agolaClientSecret"`
}

func (gitSource *CreateGitSourceRequestDto) IsValid() error {
	if len(gitSource.Name) == 0 {
		return errors.New("name is empty")
	}

	err := gitSource.GitType.IsValid()
	if err != nil {
		return err

	}

	if gitSource.GitAPIURL != nil && len(*gitSource.GitAPIURL) > 0 {
		_, err = url.ParseRequestURI(*gitSource.GitAPIURL)
		if err != nil {
			return errors.New("gitApiUrl is not valid")
		}
	} else if gitSource.GitType == types.Gitea {
		return errors.New("gitApiUrl is nil")
	}

	if len(gitSource.GitClientID) == 0 {
		return errors.New("gitClientId is empty")
	}

	if len(gitSource.GitClientSecret) == 0 {
		return errors.New("gitSecret is empty")
	}

	if gitSource.AgolaRemoteSourceName == nil {
		if gitSource.AgolaClientID == nil || gitSource.AgolaClientSecret == nil {
			return errors.New("agolaRemoteSource or oauth2 application must be specified")
		}
	}

	return nil
}
