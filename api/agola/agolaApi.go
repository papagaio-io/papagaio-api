package agola

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
)

type AgolaApiInterface interface {
	CheckOrganizationExists(agolaOrganizationRef string) bool
	CheckProjectExists(agolaOrganizationRef string, projectName string) (bool, string)
	CreateOrganization(name string, visibility dto.VisibilityType) (string, error)
	DeleteOrganization(name string, agolaUserToken string) error
	CreateProject(projectName string, organization *model.Organization, remoteSourceName string, agolaUserToken string) (string, error)
	DeleteProject(organizationName string, projectname string, agolaUserToken string) error
	GetRemoteSources() (*[]RemoteSourcesDto, error)
	AddOrUpdateOrganizationMember(agolaOrganizationRef string, agolaUserRef string, role string) error
	RemoveOrganizationMember(agolaOrganizationRef string, agolaUserRef string) error
	GetOrganizationMembers(agolaOrganizationRef string) (*OrganizationMembersResponseDto, error)
	ArchiveProject(agolaOrganizationRef string, projectName string) error
	UnarchiveProject(agolaOrganizationRef string, projectName string) error
	GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]RunDto, error)
	GetRun(runID string) (*RunDto, error)
	GetTask(runID string, taskID string) (*TaskDto, error)
	GetLogs(runID string, taskID string, step int) (string, error)
}

type AgolaApi struct{}

func (agolaApi *AgolaApi) CheckOrganizationExists(agolaOrganizationRef string) bool {
	client := &http.Client{}
	URLApi := getOrganizationUrl(agolaOrganizationRef)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	return err == nil && api.IsResponseOK(resp.StatusCode)
}

func (agolaApi *AgolaApi) CheckProjectExists(agolaOrganizationRef string, projectName string) (bool, string) {
	client := &http.Client{}
	URLApi := getProjectUrl(agolaOrganizationRef, projectName)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	var projectID string
	projectExists := err == nil && api.IsResponseOK(resp.StatusCode)
	if projectExists {
		body, _ := ioutil.ReadAll(resp.Body)
		var jsonResponse CreateProjectResponseDto //TODO make different struct with ID only
		json.Unmarshal(body, &jsonResponse)
		projectID = jsonResponse.ID
	}

	return projectExists, projectID
}

func (agolaApi *AgolaApi) CreateOrganization(name string, visibility dto.VisibilityType) (string, error) {
	client := &http.Client{}
	URLApi := getOrgUrl()
	reqBody := strings.NewReader(`{"name": "` + name + `", "visibility": "` + string(visibility) + `"}`)
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

func (agolaApi *AgolaApi) DeleteOrganization(name string, agolaUserToken string) error {
	client := &http.Client{}
	URLApi := getOrganizationUrl(name)
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

func (agolaApi *AgolaApi) CreateProject(projectName string, organization *model.Organization, remoteSourceName string, agolaUserToken string) (string, error) {
	log.Println("CreateProject start")

	if exists, projectID := agolaApi.CheckProjectExists(organization.Name, projectName); exists {
		log.Println("project already exists with ID:", projectID)
		return projectID, nil
	}

	client := &http.Client{}
	URLApi := getCreateProjectUrl()

	projectRequest := &CreateProjectRequestDto{
		Name:             projectName,
		ParentRef:        "org/" + organization.Name,
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

	fmt.Println("jsonResponse:", jsonResponse)

	return jsonResponse.ID, err
}

func (agolaApi *AgolaApi) DeleteProject(organizationName string, projectname string, agolaUserToken string) error {
	log.Println("DeleteProject start")

	client := &http.Client{}
	URLApi := getProjectUrl(organizationName, projectname)
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

func (agolaApi *AgolaApi) GetRemoteSources() (*[]RemoteSourcesDto, error) {
	client := &http.Client{}
	URLApi := getRemoteSourcesUrl()
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse []RemoteSourcesDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, err
}

func (agolaApi *AgolaApi) AddOrUpdateOrganizationMember(agolaOrganizationRef string, agolaUserRef string, role string) error {
	log.Println("AddOrUpdateOrganizationMember start")

	var err error
	client := &http.Client{}
	URLApi := getAddOrgMemberUrl(agolaOrganizationRef, agolaUserRef)
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

func (agolaApi *AgolaApi) RemoveOrganizationMember(agolaOrganizationRef string, agolaUserRef string) error {
	var err error
	client := &http.Client{}
	URLApi := getAddOrgMemberUrl(agolaOrganizationRef, agolaUserRef)
	fmt.Println("url ", URLApi)
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

func (agolaApi *AgolaApi) GetOrganizationMembers(agolaOrganizationRef string) (*OrganizationMembersResponseDto, error) {
	log.Println("GetOrganizationMembers start")

	client := &http.Client{}
	URLApi := getOrganizationMembersUrl(agolaOrganizationRef)
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
func (agolaApi *AgolaApi) ArchiveProject(agolaOrganizationRef string, projectName string) error {
	log.Println("ArchiveProject:", agolaOrganizationRef, projectName)

	return nil
}

//TODO after Agola Issue
func (agolaApi *AgolaApi) UnarchiveProject(agolaOrganizationRef string, projectName string) error {
	log.Println("UnarchiveProject:", agolaOrganizationRef, projectName)

	return nil
}

func (agolaApi *AgolaApi) GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]RunDto, error) {
	log.Println("GetRuns start")

	client := &http.Client{}
	URLApi := getRunsListUrl(projectRef, lastRun, phase, startRunID, limit, asc)
	fmt.Println("GetRuns project:", projectRef)
	fmt.Println("GetRuns url:", URLApi)
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
