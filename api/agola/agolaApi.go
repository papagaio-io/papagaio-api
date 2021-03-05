package agola

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/model"
)

//var agolaHost string = "https://agola.sorintdev.it" //TODO inserire nel config

//TODO
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

//TODO
func CreateProject(projectName string, organization *model.Organization) (string, error) {
	var err error
	var agolaProjectRef string
	return agolaProjectRef, err
}

//TODO
func DeleteProject(agolaProjectRef string, organization *model.Organization) error {
	var err error
	return err
}

//TODO
func GetRemoteSources() []string {
	var remoteSources []string
	return remoteSources
}

//TODO
func AddOrganizationMember(agolaOrganizationRef string, agolaUserRef string, role string) error {
	var err error
	return err
}

//TODO
func RemoveOrganizationMember(agolaOrganizationRef string, agolaUserRef string) error {
	var err error
	return err
}
