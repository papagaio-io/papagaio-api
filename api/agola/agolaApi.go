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

func CheckOrganizationExists(agolaOrganizationRef string) bool {
	client := &http.Client{}
	URLApi := getOrganizationUrl(agolaOrganizationRef)
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	return err == nil && api.IsResponseOK(resp.StatusCode)
}

func CreateOrganization(name string, visibility dto.VisibilityType) (string, error) {
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

func DeleteOrganization(name string, agolaUserToken string) error {
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

func CreateProject(projectName string, organization *model.Organization, remoteSourceName string, agolaUserToken string) (string, error) {
	log.Println("CreateProject start")

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

func DeleteProject(organizationName string, projectname string, agolaUserToken string) error {
	client := &http.Client{}
	URLApi := getDeleteProjectUrl(organizationName, projectname)
	req, err := http.NewRequest("DELETE", URLApi, nil)
	req.Header.Add("Authorization", "token "+agolaUserToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if !api.IsResponseOK(resp.StatusCode) {
		respMessage, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(respMessage))
	}

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

func AddOrUpdateOrganizationMember(agolaOrganizationRef string, agolaUserRef string, role string) error {
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

func RemoveOrganizationMember(agolaOrganizationRef string, agolaUserRef string) error {
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

func GetOrganizationMembers(agolaOrganizationRef string) (*OrganizationMembersResponseDto, error) {
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
func ArchiveProject(agolaOrganizationRef string, projectName string) error {
	log.Println("ArchiveProject:", agolaOrganizationRef, projectName)

	return nil
}

//TODO after Agola Issue
func UnarchiveProject(agolaOrganizationRef string, projectName string) error {
	log.Println("UnarchiveProject:", agolaOrganizationRef, projectName)

	return nil
}

func GetRuns(projectRef string, lastRun bool, phase string, startRunID *string, limit uint, asc bool) (*[]RunDto, error) {
	log.Println("GetRuns start")

	client := &http.Client{}
	URLApi := getRunsListPath(projectRef, lastRun, phase, startRunID, limit, asc)
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
