package main

import (
	catsRepo "cat-agency/internal/cats/repository"
	catsServer "cat-agency/internal/cats/server"
	catsService "cat-agency/internal/cats/service"
	"database/sql"

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
	db := openSQLConnection()

	cRepo := catsRepo.New(db)
	cService := catsService.New(cRepo)
	catsServer.InitRoutes(engine, cService)

	mRepo := missionsRepo.New(db)
	mService := missionsService.New(mRepo)
	missionsServer.InitRoutes(engine, mService)

	engine.Run(":8080")
}

func openSQLConnection() *sql.DB{
	db, err := sql.Open("sqlite3", "./../cat-agency")
	if err != nil {
		panic(err)
	}
	return db
}
