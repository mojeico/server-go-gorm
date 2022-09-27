package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=safety_service.go -destination=mocks/mock_safety_service.go

type SafetyService interface {
	GetSafetyById(companyId string, queries *helper.OneQueryParams) (models.Safety, error)
	GetAllSafeties(queries *helper.GetAllQueryParams, companyId int) ([]models.Safety, error)
	GetAllSafetiesByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	SearchSafety(searchText, offSet, limit string, companyId int) ([]models.Safety, error)
	GetAllSafetyByType(safetyType, offSet, limit string, companyId int) ([]models.Safety, error)
	GetAllSafetiesForReport(companyId int) ([]models.Safety, error)
}

type safetyService struct {
	repo     repository.SafetyRepository
	packages *initpck.Packages
}

func NewSafetyService(repo repository.SafetyRepository, packages *initpck.Packages) SafetyService {
	return &safetyService{repo: repo, packages: packages}
}

func (service *safetyService) GetAllSafetiesByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	safeties, err := service.repo.GetAllSafetiesByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedSafeties := service.packages.Sort.SortDataByStructField(orderBy, orderDir, safeties)
		return sortedSafeties, nil
	}

	return safeties, nil
}

func (service *safetyService) GetSafetyById(companyId string, queries *helper.OneQueryParams) (models.Safety, error) {

	settlement, err := service.repo.GetSafetyById(companyId, queries)

	if err != nil {
		return models.Safety{}, err
	}

	return settlement, nil
}

func (service *safetyService) GetAllSafeties(queries *helper.GetAllQueryParams, companyId int) ([]models.Safety, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	safeties, err := service.repo.GetAllSafeties(queries, companyId)

	if err != nil {
		return nil, err
	}

	return safeties, nil
}

func (service *safetyService) SearchSafety(searchText, offSet, limit string, companyId int) ([]models.Safety, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	safeties, err := service.repo.SearchSafety(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return safeties, nil

}

func (service *safetyService) GetAllSafetyByType(safetyType, offSet, limit string, companyId int) ([]models.Safety, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	safeties, err := service.repo.GetAllSafetyByType(safetyType, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return safeties, nil

}

func (service *safetyService) GetAllSafetiesForReport(companyId int) ([]models.Safety, error) {

	return service.repo.GetAllSafetiesForReport(companyId)

}
