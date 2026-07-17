package main

import (
	"core_service/internal/app"
)

// @title SkillTracker API
// @version 1.0
// @description REST API для SkillTracker.
//
// @contact.name Mandzhonson
//
// @host localhost:8080
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
