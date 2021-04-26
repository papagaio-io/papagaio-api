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

func CheckOrganizationExists(organization *model.Organization) (bool, string) {
	client := &http.Client{}
	URLApi := getOrganizationUrl(organization.AgolaOrganizationRef)
	fmt.Println("CheckOrganizationExists url:", URLApi)
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

func CheckProjectExists(organization *model.Organization, projectName string) (bool, string) {
	client := &http.Client{}
	URLApi := getProjectUrl(organization.AgolaOrganizationRef, projectName)
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

func CreateOrganization(organization *model.Organization, visibility dto.VisibilityType) (string, error) {
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

func DeleteOrganization(organization *model.Organization, agolaUserToken string) error {
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

func CreateProject(projectName string, organization *model.Organization, remoteSourceName string, agolaUserToken string) (string, error) {
	log.Println("CreateProject start")

	if exists, projectID := CheckProjectExists(organization, projectName); exists {
		log.Println("project already exists with ID:", projectID)
		return projectID, nil
	}

	client := &http.Client{}
	URLApi := getCreateProjectUrl()

	projectRequest := &CreateProjectRequestDto{
		Name:             projectName,
		ParentRef:        "org/" + organization.AgolaOrganizationRef,
		Visibility:       organization.Visibility,
		RemoteSourceName: remoteSourceName,
		RepoPath:         organization.AgolaOrganizationRef + "/" + projectName,
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

func DeleteProject(organization *model.Organization, projectname string, agolaUserToken string) error {
	log.Println("DeleteProject start")

	client := &http.Client{}
	URLApi := getProjectUrl(organization.AgolaOrganizationRef, projectname)
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

func GetRemoteSources() (*[]RemoteSourcesDto, error) {
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

func AddOrUpdateOrganizationMember(organization *model.Organization, agolaUserRef string, role string) error {
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

func RemoveOrganizationMember(organization *model.Organization, agolaUserRef string) error {
	var err error
	client := &http.Client{}
	URLApi := getAddOrgMemberUrl(organization.AgolaOrganizationRef, agolaUserRef)
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

func GetOrganizationMembers(organization *model.Organization) (*OrganizationMembersResponseDto, error) {
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
func ArchiveProject(organization *model.Organization, projectName string) error {
	log.Println("ArchiveProject:", organization.AgolaOrganizationRef, projectName)

	return nil
}

//TODO after Agola Issue
func UnarchiveProject(organization *model.Organization, projectName string) error {
	log.Println("UnarchiveProject:", organization.AgolaOrganizationRef, projectName)

	return nil
}

func GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]RunDto, error) {
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

func GetRun(runID string) (*RunDto, error) {
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

func GetTask(runID string, taskID string) (*TaskDto, error) {
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

func GetLogs(runID string, taskID string, step int) (string, error) {
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
