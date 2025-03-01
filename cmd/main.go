package main

import (
	"cat-agency/internal/cats/repository"
	"cat-agency/internal/cats/server"
	"cat-agency/internal/cats/service"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine := gin.Default()

	catsRepo := repository.New()
	catsService := service.New(catsRepo)
	server.InitRoutes(engine, catsService)

	engine.Run(":8080")
}
