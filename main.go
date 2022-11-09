package main

import "github.com/buonotti/bus-stats-api/cmd"

// swag base information

// @basePath /api/v1
// @host localhost:8080/api/v1

// @title bus-stats api
// @version 1.0
// @description The backend api for the bus-stats project

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.bearer ApiKeyAuth
// @in header
// @name api_key

func main() {
	cmd.Execute()
}
