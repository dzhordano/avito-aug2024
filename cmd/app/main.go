package main

import "github.com/dzhordano/avito-bootcamp2024/internal/app"

//	@title	Test API

//	@host						localhost:8080
//	@BasePath					/api/
//	@securityDefinitions.apikey	ClientsAuth
//	@in							header
//	@name						Authorization
//	@securityDefinitions.apikey	ModeratorsAuth
//	@in							header
//	@name						Authorization

// Main starts application through Run().
// Must specify config file path if not specified in env.
// Use -c flag to specify config file path.
func main() {
	app.Run()
}
