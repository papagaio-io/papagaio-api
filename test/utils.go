package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"wecode.sorint.it/opensource/papagaio-api/dto"
)

func ParseBody(resp *http.Response, dto interface{}) {
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(data)), dto)
}

func SortProjectsDto(projects []dto.ProjectDto) []dto.ProjectDto {
	sort.SliceStable(projects[:], func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	return projects
}

func SortBranchesDto(branches []dto.BranchDto) []dto.BranchDto {
	sort.SliceStable(branches[:], func(i, j int) bool {
		return branches[i].Name < branches[j].Name
	})
	return branches
}
