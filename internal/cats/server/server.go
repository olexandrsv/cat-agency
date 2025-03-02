package server

import (
	"cat-agency/internal/cats/service"
	"cat-agency/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	e       *gin.Engine
	service service.CatsService
}

func newServer(e *gin.Engine, service service.CatsService) *server {
	return &server{
		e:       e,
		service: service,
	}
}

func InitRoutes(e *gin.Engine, service service.CatsService) {
	s := newServer(e, service)
	e.POST("/cats", s.createCat)
	e.PUT("/cats/:id", s.updateCat)
	e.GET("/cats", s.getCats)
	e.GET("/cats/:id", s.getCat)
	e.DELETE("/cats/:id", s.deleteCat)
}

func (s *server) createCat(c *gin.Context) {
	v := common.NewValidator(c)
	experience := v.GetInt("experience")
	breed := v.GetString("breed")
	salary := v.GetFloat64("salary")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	err := s.service.CreateCat(experience, breed, salary)
	if err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func (s *server) updateCat(c *gin.Context) {
	v := common.NewValidator(c)
	id := v.GetIntFromURL("id")
	salary := v.GetFloat64("salary")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	err := s.service.UpdateCat(id, salary)
	if err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(200, "")
}

func (s *server) getCats(c *gin.Context) {
	cats, err := s.service.GetCats()
	if err != nil {
		common.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, cats)
}

func (s *server) getCat(c *gin.Context) {
	v := common.NewValidator(c)
	id := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	cat, err := s.service.GetCat(id)
	if err != nil {
		common.WriteError(c, v.Err())
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (s *server) deleteCat(c *gin.Context) {
	v := common.NewValidator(c)
	id := v.GetIntFromURL("id")

	if v.Err() != nil {
		common.WriteError(c, v.Err())
		return
	}

	err := s.service.DeleteCat(id)
	if err != nil {
		common.WriteError(c, v.Err())
		return
	}

	c.JSON(http.StatusOK, "")
}
