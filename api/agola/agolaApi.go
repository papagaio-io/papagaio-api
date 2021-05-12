package agola

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type AgolaApiInterface interface {
	CheckOrganizationExists(organization *model.Organization) (bool, string)
	CheckProjectExists(organization *model.Organization, projectName string) (bool, string)
	CreateOrganization(organization *model.Organization, visibility types.VisibilityType) (string, error)
	DeleteOrganization(organization *model.Organization, agolaUserToken string) error
	CreateProject(projectName string, agolaProjectRef string, organization *model.Organization, remoteSourceName string, agolaUserToken string) (string, error)
	DeleteProject(organization *model.Organization, agolaProjectRef string, agolaUserToken string) error
	AddOrUpdateOrganizationMember(organization *model.Organization, agolaUserRef string, role string) error
	RemoveOrganizationMember(organization *model.Organization, agolaUserRef string) error
	GetOrganizationMembers(organization *model.Organization) (*OrganizationMembersResponseDto, error)
	ArchiveProject(organization *model.Organization, agolaProjectRef string) error
	UnarchiveProject(organization *model.Organization, agolaProjectRef string) error
	GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]RunDto, error)
	GetRun(runID string) (*RunDto, error)
	GetTask(runID string, taskID string) (*TaskDto, error)
	GetLogs(runID string, taskID string, step int) (string, error)
	GetRemoteSource(agolaRemoteSource string) (*RemoteSourceDto, error)
	GetUsers() (*[]UserDto, error)
}

type AgolaApi struct{}

func (agolaApi *AgolaApi) CheckOrganizationExists(organization *model.Organization) (bool, string) {
	client := &http.Client{}
	URLApi := getOrganizationUrl(organization.AgolaOrganizationRef)

	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	var organizationID string
	organizationExists := err == nil && api.IsResponseOK(resp.StatusCode)
	if organizationExists {
		body, _ := ioutil.ReadAll(resp.Body)
		var jsonResponse AgolaCreateORGDto
		json.Unmarshal(body, &jsonResponse)
		organizationID = jsonResponse.ID
	}

	return organizationExists, organizationID
}

func (agolaApi *AgolaApi) CheckProjectExists(organization *model.Organization, agolaProjectRef string) (bool, string) {
	log.Println("CheckProjectExists start")

	client := &http.Client{}
	URLApi := getProjectUrl(organization.AgolaOrganizationRef, agolaProjectRef)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	var projectID string
	projectExists := err == nil && api.IsResponseOK(resp.StatusCode)
	if projectExists {
		body, _ := ioutil.ReadAll(resp.Body)
		var jsonResponse CreateProjectResponseDto
		json.Unmarshal(body, &jsonResponse)
		projectID = jsonResponse.ID
	}

	return projectExists, projectID
}

func (agolaApi *AgolaApi) CreateOrganization(organization *model.Organization, visibility types.VisibilityType) (string, error) {
	client := &http.Client{}
	URLApi := getOrgUrl()
	reqBody := strings.NewReader(`{"name": "` + organization.AgolaOrganizationRef + `", "visibility": "` + string(visibility) + `"}`)
	req, err := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse AgolaCreateORGDto
	json.Unmarshal(body, &jsonResponse)

	return jsonResponse.ID, err
}

func (agolaApi *AgolaApi) DeleteOrganization(organization *model.Organization, agolaUserToken string) error {
	client := &http.Client{}
	URLApi := getOrganizationUrl(organization.AgolaOrganizationRef)
	req, err := http.NewRequest("DELETE", URLApi, nil)
	req.Header.Add("Authorization", "token "+agolaUserToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	return nil
}

func (agolaApi *AgolaApi) CreateProject(projectName string, agolaProjectRef string, organization *model.Organization, remoteSourceName string, agolaUserToken string) (string, error) {
	log.Println("CreateProject start")

	if exists, projectID := agolaApi.CheckProjectExists(organization, agolaProjectRef); exists {
		log.Println("project already exists with ID:", projectID)
		return projectID, nil
	}

	client := &http.Client{}
	URLApi := getCreateProjectUrl()

	projectRequest := &CreateProjectRequestDto{
		Name:             agolaProjectRef,
		ParentRef:        "org/" + organization.AgolaOrganizationRef,
		Visibility:       organization.Visibility,
		RemoteSourceName: remoteSourceName,
		RepoPath:         organization.Name + "/" + projectName,
	}
	data, _ := json.Marshal(projectRequest)
	reqBody := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("Authorization", "token "+agolaUserToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse CreateProjectResponseDto
	json.Unmarshal(body, &jsonResponse)

	return jsonResponse.ID, err
}

func (agolaApi *AgolaApi) DeleteProject(organization *model.Organization, agolaProjectRef string, agolaUserToken string) error {
	log.Println("DeleteProject start")

	client := &http.Client{}
	URLApi := getProjectUrl(organization.AgolaOrganizationRef, agolaProjectRef)
	req, err := http.NewRequest("DELETE", URLApi, nil)
	req.Header.Add("Authorization", "token "+agolaUserToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	log.Println("DeleteProject end")

	return err
}

func (agolaApi *AgolaApi) AddOrUpdateOrganizationMember(organization *model.Organization, agolaUserRef string, role string) error {
	log.Println("AddOrUpdateOrganizationMember start")

	var err error
	client := &http.Client{}
	URLApi := getAddOrgMemberUrl(organization.AgolaOrganizationRef, agolaUserRef)
	reqBody := strings.NewReader(`{"role": "` + role + `"}`)
	req, err := http.NewRequest("PUT", URLApi, reqBody)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	log.Println("AddOrUpdateOrganizationMember end")

	return err
}

func (agolaApi *AgolaApi) RemoveOrganizationMember(organization *model.Organization, agolaUserRef string) error {
	var err error
	client := &http.Client{}
	URLApi := getAddOrgMemberUrl(organization.AgolaOrganizationRef, agolaUserRef)

	reqBody := strings.NewReader(`{}`)
	req, err := http.NewRequest("DELETE", URLApi, reqBody)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	return err
}

func (agolaApi *AgolaApi) GetOrganizationMembers(organization *model.Organization) (*OrganizationMembersResponseDto, error) {
	log.Println("GetOrganizationMembers start")

	client := &http.Client{}
	URLApi := getOrganizationMembersUrl(organization.AgolaOrganizationRef)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse OrganizationMembersResponseDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, err
}

//TODO after Agola Issue
func (agolaApi *AgolaApi) ArchiveProject(organization *model.Organization, projectName string) error {
	log.Println("ArchiveProject:", organization.AgolaOrganizationRef, projectName)

	return nil
}

//TODO after Agola Issue
func (agolaApi *AgolaApi) UnarchiveProject(organization *model.Organization, projectName string) error {
	log.Println("UnarchiveProject:", organization.AgolaOrganizationRef, projectName)

	return nil
}

func (agolaApi *AgolaApi) GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]RunDto, error) {
	log.Println("GetRuns start")

	client := &http.Client{}
	URLApi := getRunsListUrl(projectRef, lastRun, phase, startRunID, limit, asc)

	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse []RunDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, err
}

func (agolaApi *AgolaApi) GetRun(runID string) (*RunDto, error) {
	log.Println("GetRuns start")

	client := &http.Client{}
	URLApi := getRunUrl(runID)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse RunDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, err
}

func (agolaApi *AgolaApi) GetTask(runID string, taskID string) (*TaskDto, error) {
	log.Println("GetRuns start")

	client := &http.Client{}
	URLApi := getTaskUrl(runID, taskID)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse TaskDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, err
}

func (agolaApi *AgolaApi) GetLogs(runID string, taskID string, step int) (string, error) {
	log.Println("GetRuns start")

	client := &http.Client{}
	URLApi := getLogsUrl(runID, taskID, step)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(respMessage))
	}

	logs, _ := ioutil.ReadAll(resp.Body)

	return string(logs), err
}

func (agolaApi *AgolaApi) GetRemoteSource(agolaRemoteSource string) (*RemoteSourceDto, error) {
	log.Println("GetRemoteSource start")

	client := &http.Client{}
	URLApi := getRemoteSourceUrl(agolaRemoteSource)

	req, _ := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse RemoteSourceDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, nil
}

func (agolaApi *AgolaApi) GetUsers() (*[]UserDto, error) {
	log.Println("GetRemoteSource start")

	client := &http.Client{}
	URLApi := getUsersUrl()

	req, _ := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse []UserDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, nil
}
