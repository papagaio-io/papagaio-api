package agola

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type AgolaApiInterface interface {
	CheckOrganizationExists(organization *model.Organization) (bool, string)
	CheckProjectExists(organization *model.Organization, projectName string) (bool, string)
	CreateOrganization(organization *model.Organization, visibility types.VisibilityType) (string, error)
	DeleteOrganization(organization *model.Organization, user *model.User) error
	CreateProject(projectName string, agolaProjectRef string, organization *model.Organization, remoteSourceName string, user *model.User) (string, error)
	DeleteProject(organization *model.Organization, agolaProjectRef string, user *model.User) error
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
	GetOrganizations() (*[]OrganizationDto, error)

	CreateUserToken(user *model.User) error
	GetRemoteSources() (*[]RemoteSourceDto, error)
	CreateRemoteSource(remoteSourceName string, gitType string, apiUrl string, oauth2ClientId string, oauth2ClientSecret string) error
}

type AgolaApi struct {
	Db repository.Database
}

const baseTokenName string = "papagaioToken"

func (agolaApi *AgolaApi) GetOrganizations() (*[]OrganizationDto, error) {
	client := agolaApi.getClient(nil, true)
	URLApi := getOrganizationsUrl()

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse []OrganizationDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, err
}

func (agolaApi *AgolaApi) CheckOrganizationExists(organization *model.Organization) (bool, string) {
	client := agolaApi.getClient(nil, true)
	URLApi := getOrganizationUrl(organization.AgolaOrganizationRef)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)
	if err != nil {
		return false, ""
	}

	defer resp.Body.Close()

	var organizationID string
	organizationExists := api.IsResponseOK(resp.StatusCode)
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

	client := agolaApi.getClient(nil, true)
	URLApi := getProjectUrl(organization.AgolaOrganizationRef, agolaProjectRef)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	var projectID string
	projectExists := api.IsResponseOK(resp.StatusCode)
	if projectExists {
		body, _ := ioutil.ReadAll(resp.Body)
		var jsonResponse CreateProjectResponseDto
		json.Unmarshal(body, &jsonResponse)
		projectID = jsonResponse.ID
	}

	return projectExists, projectID
}

func (agolaApi *AgolaApi) CreateOrganization(organization *model.Organization, visibility types.VisibilityType) (string, error) {
	client := agolaApi.getClient(nil, true)
	URLApi := getOrgUrl()
	reqBody := strings.NewReader(`{"name": "` + organization.AgolaOrganizationRef + `", "visibility": "` + string(visibility) + `"}`)
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
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

func (agolaApi *AgolaApi) DeleteOrganization(organization *model.Organization, user *model.User) error {
	client := agolaApi.getClient(user, false)
	URLApi := getOrganizationUrl(organization.AgolaOrganizationRef)
	req, _ := http.NewRequest("DELETE", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
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

func (agolaApi *AgolaApi) CreateProject(projectName string, agolaProjectRef string, organization *model.Organization, remoteSourceName string, user *model.User) (string, error) {
	log.Println("CreateProject start")

	if exists, projectID := agolaApi.CheckProjectExists(organization, agolaProjectRef); exists {
		log.Println("project already exists with ID:", projectID)
		return projectID, nil
	}

	client := agolaApi.getClient(user, false)
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

	req, _ := http.NewRequest("POST", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
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

func (agolaApi *AgolaApi) DeleteProject(organization *model.Organization, agolaProjectRef string, user *model.User) error {
	log.Println("DeleteProject start")

	client := agolaApi.getClient(user, false)
	URLApi := getProjectUrl(organization.AgolaOrganizationRef, agolaProjectRef)
	req, _ := http.NewRequest("DELETE", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
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

	log.Println("AddOrUpdateOrganizationMember", agolaUserRef, "for", organization.Name, "with role:", role)

	var err error
	client := agolaApi.getClient(nil, true)
	URLApi := getAddOrgMemberUrl(organization.AgolaOrganizationRef, agolaUserRef)
	reqBody := strings.NewReader(`{"role": "` + role + `"}`)
	req, _ := http.NewRequest("PUT", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	log.Println("AddOrUpdateOrganizationMember end")

	return err
}

func (agolaApi *AgolaApi) RemoveOrganizationMember(organization *model.Organization, agolaUserRef string) error {
	log.Println("RemoveOrganizationMember", organization.Name, "with agolaUserRef", agolaUserRef)

	var err error
	client := agolaApi.getClient(nil, true)
	URLApi := getAddOrgMemberUrl(organization.AgolaOrganizationRef, agolaUserRef)

	reqBody := strings.NewReader(`{}`)
	req, _ := http.NewRequest("DELETE", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("RemoveOrganizationMember StatusCode:", resp.StatusCode)

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		log.Println("RemoveOrganizationMember respMessage:", string(respMessage))
		return errors.New(string(respMessage))
	}

	return err
}

func (agolaApi *AgolaApi) GetOrganizationMembers(organization *model.Organization) (*OrganizationMembersResponseDto, error) {
	log.Println("GetOrganizationMembers start")

	client := agolaApi.getClient(nil, true)
	URLApi := getOrganizationMembersUrl(organization.AgolaOrganizationRef)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
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

	client := agolaApi.getClient(nil, true)
	URLApi := getRunsListUrl(projectRef, lastRun, phase, startRunID, limit, asc)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
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

	client := agolaApi.getClient(nil, true)
	URLApi := getRunUrl(runID)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
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

	client := agolaApi.getClient(nil, true)
	URLApi := getTaskUrl(runID, taskID)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
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

	client := agolaApi.getClient(nil, true)
	URLApi := getLogsUrl(runID, taskID, step)
	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
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

	client := agolaApi.getClient(nil, true)
	URLApi := getRemoteSourceUrl(agolaRemoteSource)

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
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

	client := agolaApi.getClient(nil, true)
	URLApi := getUsersUrl()

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
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

func (agolaApi *AgolaApi) CreateUserToken(user *model.User) error {
	if user == nil || user.AgolaUserRef == nil {
		log.Println("CreateUserToken error user nil")
		return errors.New("user nil error")
	}

	log.Println("RefreshAgolaUserToken user", *user.AgolaUserRef)

	if user.UserID == nil {
		log.Println("UserID is nil")
		return errors.New("UsersID is nil")
	}
	if user.AgolaUserRef == nil {
		log.Println("AgolaUserRef is nil")
		return errors.New("AgolaUserRef is nil")
	}

	tokenName := baseTokenName + "-" + fmt.Sprint(time.Now().Unix())
	user.AgolaTokenName = &tokenName

	client := &http.Client{}
	URLApi := getCreateTokenUrl(*user.AgolaUserRef)

	tokenRequest := &TokenRequestDto{
		TokenName: *user.AgolaTokenName,
	}
	data, _ := json.Marshal(tokenRequest)
	reqBody := strings.NewReader(string(data))

	req, _ := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Set("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse TokenResponseDto
	json.Unmarshal(body, &jsonResponse)

	user.AgolaToken = &jsonResponse.Token
	agolaApi.Db.SaveUser(user)

	return nil
}

func (agolaApi *AgolaApi) GetRemoteSources() (*[]RemoteSourceDto, error) {
	log.Println("GetRemoteSources start")

	client := agolaApi.getClient(nil, true)
	URLApi := getRemoteSourcesUrl()

	req, _ := http.NewRequest("GET", URLApi, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(respMessage))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var jsonResponse []RemoteSourceDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse, nil
}

func (agolaApi *AgolaApi) CreateRemoteSource(remoteSourceName string, gitType string, apiUrl string, oauth2ClientId string, oauth2ClientSecret string) error {
	log.Println("CreateRemoteSource start")

	client := agolaApi.getClient(nil, true)
	URLApi := getRemoteSourcesUrl()

	projectRequest := &CreateRemoteSourceRequestDto{
		Name:                remoteSourceName,
		APIURL:              apiUrl,
		Type:                gitType,
		AuthType:            "oauth2",
		SkipSSHHostKeyCheck: true,
		SkipVerify:          false,
		Oauth2ClientID:      oauth2ClientId,
		Oauth2ClientSecret:  oauth2ClientSecret,
	}
	data, _ := json.Marshal(projectRequest)
	reqBody := strings.NewReader(string(data))

	req, _ := http.NewRequest("POST", URLApi, reqBody)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

	return nil
}

///////////////

func (agolaApi *AgolaApi) getClient(user *model.User, isAdminUser bool) *httpClient {
	client := &httpClient{c: &http.Client{}, user: user, agolaApi: agolaApi, isAdminUser: isAdminUser}

	return client
}

type httpClient struct {
	c           *http.Client
	user        *model.User
	agolaApi    *AgolaApi
	isAdminUser bool
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	if c.isAdminUser {
		req.Header.Set("Authorization", config.Config.Agola.AdminToken)
		return c.c.Do(req)
	}

	fmt.Println("user before:", *c.user.AgolaTokenName, ",", *c.user.AgolaToken)

	req.Header.Set("Authorization", "token "+*c.user.AgolaToken)
	response, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("statucCode:", response.StatusCode)

	if response.StatusCode == 401 {
		err = c.agolaApi.CreateUserToken(c.user)
		if err != nil {
			log.Println("error in agola CreateUserToken:", err)
			return nil, err
		}

		req.Header.Set("Authorization", "token "+*c.user.AgolaToken)
		response, err = c.c.Do(req)
	}

	fmt.Println("user after:", *c.user.AgolaTokenName, ",", *c.user.AgolaToken)

	return response, err
}
