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
	e.POST("/missions", s.createMission)
	e.DELETE("/missions/:id", s.deleteMission)
	e.PUT("/missions/:id/cat", s.assignMission)
	e.PUT("/missions/:id/completed", s.updateMission)
	e.PUT("/targets/:id/completed", s.updateTarget)
	e.DELETE("/targets/:id", s.deleteTarget)
	e.POST("/missions/:mission_id/targets", s.createTarget)
	e.GET("/missions", s.getMissions)
	e.GET("/missions/:id", s.getMission)
	e.PUT("/missions/:id/targets/:target_id/notes/:note_id", s.updateNote)
}

func (s *server) createMission(c *gin.Context){
	var targets []models.Target
	if err := c.BindJSON(&targets); err != nil {
		common.WriteError(c, common.NewJSONParseError(err))
		return
	}

	if err := s.service.CreateMission(targets); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func (s *server) deleteMission(c *gin.Context){
	v := common.NewValidator(c)
	id := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	if err := s.service.DeleteMission(id); err != nil {
		common.WriteError(c, err)
	}

	c.JSON(http.StatusOK, "")
}

func (s *server) assignMission(c *gin.Context){
	v := common.NewValidator(c)
	missionID := v.GetIntFromURL("id")
	catID := v.GetInt("cat_id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}
	
	if err := s.service.AssignMission(missionID, catID); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func (s *server) updateMission(c *gin.Context){
	v := common.NewValidator(c)
	missionID := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	if err := s.service.UpdateMission(missionID); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, "")
}

func (s *server) updateTarget(c *gin.Context){
	v := common.NewValidator(c)
	missionID := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	if err := s.service.UpdateTarget(missionID); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, "")
}

func (s *server) deleteTarget(c *gin.Context){
	v := common.NewValidator(c)
	targetID := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	if err := s.service.DeleteTarget(targetID); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, "")
}

func (s *server) createTarget(c *gin.Context){
	v := common.NewValidator(c)
	missionID := v.GetIntFromURL("mission_id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	var target models.Target
	if err := c.BindJSON(&target); err != nil {
		common.WriteError(c, common.NewJSONParseError(err))
		return
	}

	if err := s.service.CreateTarget(missionID, target); err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, "")
}

func (s *server) getMissions(c *gin.Context){
	missions, err := s.service.GetMissions()
	if  err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, missions)
}

func (s *server) getMission(c *gin.Context){
	v := common.NewValidator(c)
	id := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	mission, err := s.service.GetMission(id)
	if err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, mission)
}

func (s *server) updateNote(c *gin.Context){
	v := common.NewValidator(c)
	missionID := v.GetIntFromURL("id")
	targetID := v.GetIntFromURL("target_id")
	noteID := v.GetIntFromURL("note_id")
	msg := v.GetString("msg")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	err := s.service.UpdateNote(missionID, targetID, noteID, msg)
	if err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, "")
}