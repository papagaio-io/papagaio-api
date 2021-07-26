# Papagaio-API

# Commands for mock build

mockgen -source repository/database.go -destination .\test\mock\mock_repository\mockDatabase.go
mockgen -source api/agola/agolaApi.go -destination .\test\mock\mock_agola\mockAgolaApi.go
mockgen -source api/git/gitea/giteaApi.go -destination .\test\mock\mock_gitea\mockGiteaApi.go
mockgen -source api/git/github/githubApi.go -destination .\test\mock\mock_github\mockGithubApi.go
mockgen -source api/git/gitlab/gitlabApi.go -destination .\test\mock\mock_gitlab\mockGitlabApi.go

# Configuration

* Add a gitSource with this command
papagaio gitsource add  
      --agola-client-id string       agola oauth2 client id
      --agola-client-secret string   agola oauth2 client secret
      --agola-remotesource string    agola remotesource name
      --gateway-url string           papagaio gateway URL(optional)
      --git-api-url string           api url
      --git-client-id string         git oauth2 client id
      --git-client-secret string     git oauth2 client secret
  -h, --help                         help for gitsource
      --name string                  gitSource name
      --token string                 token
      --type string                  git type(gitea, github, gitlab)
      --delete-remotesource          true to delete the Agola remotesource(default false)


example: papagaio gitsource add --name {gitSourceName} --type gitea --git-token {gitUserToken} --git-api-url  {gitUrl --agola-remotesource {agolaRemoteSource} --agola-token {agolaUserToken} --token {papagaioAdminToken}

* Change user role
papagaio user change-role
      --gateway-url string   papagaio gateway URL(optional)
      --id uint              user id
      --role string          user role(ADMINISTRATOR, DEVELOPER)
      -h, --help   help for change-role
      --token string         token
example: papagaio user change-role --id {userId} --role ADMINISTRATOR --token {papagaioAdminToken}

# Swagger

* Use command line "swag init" to update swag autogenerate files
* URL: /swagger/index.html or /swagger/

# Test

go test wecode.sorint.it/opensource/papagaio-api/service -v

go test -coverprofile testCover.out wecode.sorint.it/opensource/papagaio-api/service
go tool cover -html testCover.out