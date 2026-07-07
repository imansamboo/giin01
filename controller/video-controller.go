package controller

import (
	"log"

	"godlv/entity"
	"godlv/service"
	localValidator "godlv/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type VideoController interface {
	FindAll(ctx *gin.Context) []entity.Video
	Save(ctx *gin.Context) error
}

type videoController struct {
	videoService service.VideoService
}

var customValidator *validator.Validate

func New(service service.VideoService) VideoController {
	customValidator = validator.New()
	customValidator.RegisterValidation("happy", localValidator.ValidateHappy)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("happy", localValidator.ValidateHappy)
	}
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
	log.Println("[save] binding request body")
	var video entity.Video
	if err := ctx.BindJSON(&video); err != nil {
		log.Printf("[save] bind JSON failed: %v", err)
		return err
	}
	if err := customValidator.Struct(video); err != nil {
		log.Printf("[save] validation failed: %v", err)
		return err
	}
	log.Printf("[save] bound video title=%q url=%q", video.Title, video.URL)

	video, err := c.videoService.Save(video)
	if err != nil {
		log.Printf("[save] service save failed: %v", err)
		return err
	}
	log.Printf("[save] service saved video title=%q", video.Title)
	return nil
}
