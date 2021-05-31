package service

import (
	"myMod/entity"
	"myMod/repository"
)

type VideoService interface {
	Create(video entity.Video) entity.Video
	FindAll() []entity.Video
}

type videoService struct {
	videoRepository repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) VideoService {
	return &videoService{
		videoRepository: repo,
	}
}

func (service *videoService) Create(video entity.Video) entity.Video {
	service.videoRepository.Create(video)
	return video
}

func (service *videoService) FindAll() []entity.Video {
	return service.videoRepository.FindAll()
}
