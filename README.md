# Papagaio-API

# Commands for mock build
mockgen -source repository/database.go -destination .\test\mock\mock_repository\mockDatabase.go
mockgen -source api/agola/agolaApi.go -destination .\test\mock\mock_agola\mockAgolaApi.go
mockgen -source api/git/gitea/giteaApi.go -destination .\test\mock\mock_gitea\mockGiteaApi.go
mockgen -source api/git/github/githubApi.go -destination .\test\mock\mock_github\mockGithubApi.go