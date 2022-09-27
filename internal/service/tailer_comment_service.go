package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	initpck "github.com/trucktrace/pkg/main"
)

type TrailerCommentService interface {
	GetTrailerCommentsByTrailerID(trailerId string) ([]models.TrailerComment, error)
}

type trailerCommentService struct {
	repo    repository.TrailerCommentRepository
	initpck *initpck.Packages
}

func NewTrailerCommentService(repo repository.TrailerCommentRepository, initpck *initpck.Packages) TrailerCommentService {
	return &trailerCommentService{
		repo:    repo,
		initpck: initpck,
	}
}

func (service trailerCommentService) GetTrailerCommentsByTrailerID(trailerId string) ([]models.TrailerComment, error) {

	return service.repo.GetTrailerCommentsByTrailerID(trailerId)
}
