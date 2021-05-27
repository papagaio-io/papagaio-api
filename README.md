# Papagaio-API

# Commands for mock build

mockgen -source repository/database.go -destination .\test\mock\mock_repository\mockDatabase.go
mockgen -source api/agola/agolaApi.go -destination .\test\mock\mock_agola\mockAgolaApi.go
mockgen -source api/git/gitea/giteaApi.go -destination .\test\mock\mock_gitea\mockGiteaApi.go
mockgen -source api/git/github/githubApi.go -destination .\test\mock\mock_github\mockGithubApi.go

# Configuration

* Add a gitSource with this command
papagaio gitsource add  
      --agola-remotesource string   Agola remotesource
      --agola-token string          Agola token
      --gateway-url string          papagaio gateway URL(optional)
      --git-api-url string          api url
      --git-token string            git token
  -h, --help                        help for gitsource
      --name string                 gitSource name
      --token string                token
      --type string                 git type(gitea, github)

example: papagaio gitsource add --name {gitSourceName} --type gitea --git-token {gitUserToken} --git-api-url  {gitUrl --agola-remotesource {agolaRemoteSource} --agola-token {agolaUserToken} --token {papagaioAdminToken}

* Add administration users
papagaio user add
      --email string         user email
      --gateway-url string   papagaio gateway URL(optional)
  -h, --help                 help for user
      --token string         token
example: papagaio user add --email {userEmail} --token {papagaioAdminToken}

# Swagger

* Use command line "swag init" to update swaw autogenerate files
* URL: /swagger/index.html or /swagger/