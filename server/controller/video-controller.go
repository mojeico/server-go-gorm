package controller

import (
	"github.com/gin-gonic/gin"
	"myMod/entity"
	"myMod/service"
)

type VideoController interface {
	FindAll() []entity.Video
	Create(ctx *gin.Context) error
}

type controller struct {
	service service.VideoService
}

func NewVideoController(service service.VideoService) VideoController {
	return &controller{
		service: service,
	}
}

func (c *controller) FindAll() []entity.Video {
	return c.service.FindAll()
}

func (c *controller) Create(ctx *gin.Context) error {
	var video entity.Video

	if err := ctx.BindJSON(&video); err != nil {
		return err
	}
	c.service.Create(video)
	return nil
}
