package service

import (
	"log"

	"godlv/entity"
)

type VideoService interface {
	FindAll() ([]entity.Video, error)
	Save(video entity.Video) (entity.Video, error)
}

type videoService struct {
	videos []entity.Video
}

func New() VideoService {
	return &videoService{}
}

func (service *videoService) FindAll() ([]entity.Video, error) {
	return service.videos, nil
}

func (service *videoService) Save(video entity.Video) (entity.Video, error) {
	log.Printf("[save] persisting video title=%q (total before=%d)", video.Title, len(service.videos))
	service.videos = append(service.videos, video)
	log.Printf("[save] video persisted (total after=%d)", len(service.videos))
	return video, nil
}
