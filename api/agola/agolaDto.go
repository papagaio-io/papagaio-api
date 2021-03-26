package agola

import (
	"strings"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/dto"
)

type AgolaCreateORGDto struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Visibility dto.VisibilityType `json:"visibility"`
}

type RemoteSourcesDto struct {
	Name string `json:"name"`
}

type CreateProjectRequestDto struct {
	Name             string             `json:"name"`
	ParentRef        string             `json:"parent_ref"`
	Visibility       dto.VisibilityType `json:"visibility"`
	RemoteSourceName string             `json:"remote_source_name"`
	RepoPath         string             `json:"repo_path"`
}

type CreateProjectResponseDto struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Path             string             `json:"path"`
	ParentPath       string             `json:"parent_path"`
	Visibility       dto.VisibilityType `json:"visibility"`
	GlobalVisibility string             `json:"global_visibility"`
}

type OrganizationMembersResponseDto struct {
	Members []MemberDto `json:"members"`
}

type MemberDto struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Role     RoleType `json:"role"`
}

type RoleType string

const (
	Owner  RoleType = "owner"
	Member RoleType = "member"
)

type RunDto struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Counter     uint64            `json:"counter"`
	Annotations map[string]string `json:"annotations"`
	Phase       RunPhase          `json:"phase"`
	Result      RunResult         `json:"result"`
	StartTime   *time.Time        `json:"start_time"`
	EndTime     *time.Time        `json:"end_time"`
	Archived    bool              `json:"archived"`
}

type RunPhase string

const (
	RunPhaseSetupError RunPhase = "setuperror"
	RunPhaseQueued     RunPhase = "queued"
	RunPhaseCancelled  RunPhase = "cancelled"
	RunPhaseRunning    RunPhase = "running"
	RunPhaseFinished   RunPhase = "finished"
)

type RunResult string

const (
	RunResultUnknown RunResult = "unknown"
	RunResultStopped RunResult = "stopped"
	RunResultSuccess RunResult = "success"
	RunResultFailed  RunResult = "failed"
)

func (run *RunDto) IsWebhookCreationTrigger() bool {
	return strings.Compare(run.Annotations["run_creation_trigger"], "webhook") == 0
}

func (run *RunDto) GetBranchName() string {
	return run.Annotations["branch"]
}

func (run *RunDto) GetCommitSha() string {
	return run.Annotations["commit_sha"]
}
