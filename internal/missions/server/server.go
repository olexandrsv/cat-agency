package server

import (
	"cat-agency/internal/common"
	"cat-agency/internal/missions/models"
	"cat-agency/internal/missions/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	e *gin.Engine
	service service.MissionService
}

func newServer(e *gin.Engine, service service.MissionService) *server {
	return &server{
		e:       e,
		service: service,
	}
}

func InitRoutes(e *gin.Engine, service service.MissionService) {
	s := newServer(e, service)
	e.POST("/mission", s.createMission)
}

func (s *server) createMission(c *gin.Context){
	var targets []models.Target
	if err := c.BindJSON(&targets); err != nil {
		common.WriteError(c, err)
		return
	}

	if err := s.service.CreateMission(targets); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}