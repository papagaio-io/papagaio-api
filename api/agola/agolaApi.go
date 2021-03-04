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
	var agolaOrganizationID string
	var err error

	return agolaOrganizationID, err
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

func CreateUserToken(agolaUserRef string, tokenName string) (string, error) {
	client := &http.Client{}

	URLApi := getCreateTokenUrl(agolaUserRef)

	reqBody := strings.NewReader(`{"token_name": "` + tokenName + `"}`)
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

	var jsonResponse AgolaCreateTokenDto
	json.Unmarshal(body, &jsonResponse)

	return jsonResponse.Token, err
}

//TODO
func DeleteUserToken(agolaUserRef string, tokenName string) error {
	var err error
	return err
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
