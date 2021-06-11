package agola

import (
	"strings"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/types"
)

type AgolaCreateORGDto struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Visibility types.VisibilityType `json:"visibility"`
}

type RemoteSourcesDto struct {
	Name string `json:"name"`
}

type CreateProjectRequestDto struct {
	Name             string               `json:"name"`
	ParentRef        string               `json:"parent_ref"`
	Visibility       types.VisibilityType `json:"visibility"`
	RemoteSourceName string               `json:"remote_source_name"`
	RepoPath         string               `json:"repo_path"`
}

type CreateProjectResponseDto struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Path             string               `json:"path"`
	ParentPath       string               `json:"parent_path"`
	Visibility       types.VisibilityType `json:"visibility"`
	GlobalVisibility string               `json:"global_visibility"`
}

type OrganizationMembersResponseDto struct {
	Members []MemberDto `json:"members"`
}

type MemberDto struct {
	User UserDto  `json:"user"`
	Role RoleType `json:"role"`
}

type RoleType string

const (
	Owner  RoleType = "owner"
	Member RoleType = "member"
)

type RunDto struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Counter     uint64             `json:"counter"`
	Annotations map[string]string  `json:"annotations"`
	Tasks       map[string]TaskDto `json:"tasks"`
	Phase       RunPhase           `json:"phase"`
	Result      RunResult          `json:"result"`
	StartTime   time.Time          `json:"start_time"`
	EndTime     time.Time          `json:"end_time"`
	Archived    bool               `json:"archived"`
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

type TaskDto struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Status    RunTaskStatus `json:"status"`
	Steps     []RunTaskStep `json:"steps"`
	SetupStep RunTaskStep   `json:"setup_step"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
}

type RunTaskStatus string

const (
	RunTaskStatusNotStarted RunTaskStatus = "notstarted"
	RunTaskStatusSkipped    RunTaskStatus = "skipped"
	RunTaskStatusCancelled  RunTaskStatus = "cancelled"
	RunTaskStatusRunning    RunTaskStatus = "running"
	RunTaskStatusStopped    RunTaskStatus = "stopped"
	RunTaskStatusSuccess    RunTaskStatus = "success"
	RunTaskStatusFailed     RunTaskStatus = "failed"
)

type RunTaskStep struct {
	Phase      ExecutorTaskPhase `json:"phase"`
	Name       string            `json:"name"`
	LogPhase   RunTaskFetchPhase `json:"log_phase"`
	ExitStatus int               `json:"exit_status"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
}

type ExecutorTaskPhase string

const (
	ExecutorTaskPhaseNotStarted ExecutorTaskPhase = "notstarted"
	ExecutorTaskPhaseCancelled  ExecutorTaskPhase = "cancelled"
	ExecutorTaskPhaseRunning    ExecutorTaskPhase = "running"
	ExecutorTaskPhaseStopped    ExecutorTaskPhase = "stopped"
	ExecutorTaskPhaseSuccess    ExecutorTaskPhase = "success"
	ExecutorTaskPhaseFailed     ExecutorTaskPhase = "failed"
)

type RunTaskFetchPhase string

const (
	RunTaskFetchPhaseNotStarted RunTaskFetchPhase = "notstarted"
	RunTaskFetchPhaseFinished   RunTaskFetchPhase = "finished"
)

type RemoteSourceDto struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	AuthType            string `json:"auth_type"`
	RegistrationEnabled bool   `json:"registration_enabled"`
	LoginEnabled        bool   `json:"login_enabled"`
}

type UserDto struct {
	ID             string             `json:"id"`
	Username       string             `json:"username"`
	LinkedAccounts []LinkedAccountDto `json:"linked_accounts"`
}

type LinkedAccountDto struct {
	ID                  string `json:"id"`
	RemoteSourceID      string `json:"remote_source_id"`
	RemoteUserName      string `json:"remote_user_name"`
	RemoteUserAvatarURL string `json:"remote_user_avatar_url"`
}

type OrganizationDto struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

type TokenRequestDto struct {
	TokenName string `json:"token_name"`
}

type TokenResponseDto struct {
	Token string `json:"token"`
}

type CreateRemoteSourceRequestDto struct {
	Name                string `json:"name"`
	APIURL              string `json:"apiurl"`
	Type                string `json:"type"`
	AuthType            string `json:"auth_type"`
	SkipVerify          bool   `json:"skip_verify"`
	Oauth2ClientID      string `json:"oauth_2_client_id"`
	Oauth2ClientSecret  string `json:"oauth_2_client_secret"`
	SSHHostKey          string `json:"ssh_host_key"`
	SkipSSHHostKeyCheck bool   `json:"skip_ssh_host_key_check"`
	RegistrationEnabled *bool  `json:"registration_enabled"`
	LoginEnabled        *bool  `json:"login_enabled"`
}
