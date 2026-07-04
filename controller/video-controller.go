package controller

import (
	"godlv/entity"
	"godlv/service"

	"github.com/gin-gonic/gin"
)

type VideoController interface {
	FindAll(ctx *gin.Context) []entity.Video
	Save(ctx *gin.Context) error
}

type videoController struct {
	videoService service.VideoService
}

func New(service service.VideoService) VideoController {
	return &videoController{
		videoService: service,
	}
}

func (c *videoController) FindAll(ctx *gin.Context) []entity.Video {
	videos, err := c.videoService.FindAll()
	if err != nil {
		return nil
	}
	return videos
}

func (c *videoController) Save(ctx *gin.Context) error {
	var video entity.Video
	if err := ctx.ShouldBindJSON(&video); err != nil {
		return err
	}
	video, err := c.videoService.Save(video)
	if err != nil {
		return err
	}
	return nil

}
