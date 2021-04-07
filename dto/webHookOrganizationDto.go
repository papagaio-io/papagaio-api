package dto

import "strings"

type WebHookDto struct {
	Secret     string        `json:"secret"`
	Sha        string        `json:"sha"`
	Action     string        `json:"action"`
	Repository RepositoryDto `json:"repository"`
	RefType    string        `json:"ref_type"`
}

func (webHookMessage *WebHookDto) IsRepositoryCreated() bool {
	return strings.Compare(webHookMessage.Action, "created") == 0
}

func (webHookMessage *WebHookDto) IsRepositoryDeleted() bool {
	return strings.Compare(webHookMessage.Action, "deleted") == 0
}

func (webHookMessage *WebHookDto) IsPush() bool {
	return strings.Compare(webHookMessage.Action, "") == 0
}

type RepositoryDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
