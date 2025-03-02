package main

import (
	catsRepo "cat-agency/internal/cats/repository"
	catsServer "cat-agency/internal/cats/server"
	catsService "cat-agency/internal/cats/service"
	// "fmt"

	// "cat-agency/internal/missions/models"
	missionsRepo "cat-agency/internal/missions/repository"
	missionsServer "cat-agency/internal/missions/server"
	missionsService "cat-agency/internal/missions/service"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine := gin.Default()

	cRepo := catsRepo.New()
	cService := catsService.New(cRepo)
	catsServer.InitRoutes(engine, cService)

	mRepo := missionsRepo.New()
	mService := missionsService.New(mRepo)
	missionsServer.InitRoutes(engine, mService)

	engine.Run(":8080")
}
