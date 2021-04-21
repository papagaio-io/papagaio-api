# Papagaio-API

# Commands for mock build
mockgen -source repository/database.go -destination .\test\mock\mockDatabase.go
mockgen -source api/agola/agolaApi.go -destination .\test\mock\mockAgolaApi.go
mockgen -source api/git/gitea/giteaApi.go -destination .\test\mock\mockGiteaApi.go
mockgen -source api/git/github/githubApi.go -destination .\test\mock\mockGithubApi.go