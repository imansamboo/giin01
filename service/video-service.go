package service

import "godlv/entity"

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
	service.videos = append(service.videos, video)
	return video, nil
}
