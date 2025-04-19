// @title MCP API
// @version 1.0
// @description API сервер для MCP
// @BasePath /
package main

import (
	"api/internal/config"
	"api/internal/infrastructure"
	"time"
)

func main() {
	config, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}

	app := infrastructure.NewConfig(config).
		ContextTimeout(10 * time.Second).
		Logger().
		Validator()

	_ = app
}

/*
TODO :
1. Add a database connection (PostgreSQL)
2. Add repo, entities
3. Add migrations (user, organization)
4. Add usecases/actions/routes
*/
