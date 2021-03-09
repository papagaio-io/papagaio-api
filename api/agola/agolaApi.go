package agola

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

func CreateOrganization(name string, visibility string) (string, error) {
	client := &http.Client{}
	URLApi := getCreateORGUrl()
	reqBody := strings.NewReader(`{"name": "` + name + `", "visibility": "` + visibility + `"}`)
	req, err := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	if resp.StatusCode == 400 {
		return "", errors.New(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse AgolaCreateORGDto
	json.Unmarshal(body, &jsonResponse)

	return jsonResponse.ID, err
}

//TODO
func GetOrganizations() (*[]model.Organization, error) {
	var organizations *[]model.Organization
	var err error
	return organizations, err
}

func CreateProject(projectName string, organization *model.Organization, remoteSourceName string) (string, error) {
	client := &http.Client{}
	URLApi := getCreateProjectPath()

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
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	if resp.StatusCode == 400 {
		return "", errors.New(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse CreateProjectResponseDto
	json.Unmarshal(body, &jsonResponse)

	return jsonResponse.ID, err
}

//TODO
func DeleteProject(agolaProjectRef string, organization *model.Organization) error {
	var err error
	return err
}

func GetRemoteSources() *[]RemoteSourcesDto {
	client := &http.Client{}
	URLApi := getRemoteSourcesUrl()
	req, err := http.NewRequest("GET", URLApi, nil)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse []RemoteSourcesDto
	json.Unmarshal(body, &jsonResponse)

	return &jsonResponse
}

func AddOrganizationMember(agolaOrganizationRef string, agolaUserRef string, role string) error {
	var err error
	client := &http.Client{}
	URLApi := getAddOrgMemberUrl(agolaOrganizationRef, agolaUserRef)
	reqBody := strings.NewReader(`{"role": "` + role + `"}`)
	req, err := http.NewRequest("PUT", URLApi, reqBody)
	req.Header.Add("Authorization", config.Config.Agola.AdminToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode == 400 {
		return errors.New(resp.Status)
	}

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

	if err != nil {
		return err
	}

	if resp.StatusCode == 400 {
		return errors.New(resp.Status)
	}

	return err
}
