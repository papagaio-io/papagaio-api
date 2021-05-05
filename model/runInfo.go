package model

import (
	"fmt"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/types"
)

type RunInfo struct {
	ID           string          `json:"id"`
	Branch       string          `json:"branch"`
	RunStartDate time.Time       `json:"runStartDate"`
	RunEndDate   time.Time       `json:"runEndDate,omitempty"`
	Phase        types.RunPhase  `json:"phase"`
	Result       types.RunResult `json:"result"`
}

const runURL string = "%s/org/%s/projects/%s.proj/runs/%s"

func (run *RunInfo) GetURL(organization *Organization, project *Project) string {
	return fmt.Sprintf(runURL, config.Config.Agola.AgolaAddr, organization.AgolaOrganizationRef, project.AgolaProjectRef, run.ID)
}
