package main

import (
	_ "github.com/swaggo/files" // swagger embed files

	"wecode.sorint.it/opensource/papagaio-api/cmd"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/docs"
	_ "wecode.sorint.it/opensource/papagaio-api/docs" // docs is generated by Swag CLI, you have to import it.
)

// @BasePath /api
// @title papagaio-api
// @version 0.1.0

// @securitydefinitions.apiKey ApiKeyToken
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.Host = config.Config.Server.LocalHostAddress

	cmd.Execute()
}
