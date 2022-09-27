package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=trailer_service.go -destination=mocks/mock_trailer_service.go

type TrailerService interface {
	GetAllTrailerByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	GetTrailerById(queries *helper.OneQueryParams) (map[string]interface{}, error)
	GetAllTrailers(queries *helper.GetAllQueryParams, companyId int) ([]map[string]interface{}, error)
	SearchTrailers(searchText, offSet, limit string, companyId int) ([]map[string]interface{}, error)
	GetAllTrailersForReport(companyId int) ([]models.Trailer, error)
}

type trailerService struct {
	repo     repository.TrailerRepository
	packages *initpck.Packages
}

func NewTrailerService(repo repository.TrailerRepository, packages *initpck.Packages) TrailerService {
	return &trailerService{
		repo:     repo,
		packages: packages,
	}
}

func (service *trailerService) GetTrailerById(queries *helper.OneQueryParams) (map[string]interface{}, error) {

	trailer, err := service.repo.GetTrailerById(queries)

	if err != nil {
		return nil, err
	}

	return trailer, nil
}

func (service *trailerService) GetAllTrailerByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	trailers, err := service.repo.GetAllTrailerByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedTrailers := service.packages.Sort.SortDataByStructField(orderBy, orderDir, trailers)
		return sortedTrailers, nil
	}

	return trailers, nil
}

func (service *trailerService) GetAllTrailers(queries *helper.GetAllQueryParams, companyId int) ([]map[string]interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}
	trailers, err := service.repo.GetAllTrailers(queries, companyId)

	if err != nil {
		return nil, err
	}

	return trailers, nil
}

func (service *trailerService) SearchTrailers(searchText, offSet, limit string, companyId int) ([]map[string]interface{}, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	trailers, err := service.repo.SearchTrailers(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return trailers, nil
}

func (service *trailerService) GetAllTrailersForReport(companyId int) ([]models.Trailer, error) {

	return service.repo.GetAllTrailersForReport(companyId)

}
