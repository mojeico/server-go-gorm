package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=driver_service.go -destination=mocks/mock_driver_service.go

type DriverService interface {
	GetDriverById(queries *helper.OneQueryParams) (models.Driver, error)
	GetAllDrivers(queries *helper.GetAllQueryParams, companyId int) ([]models.Driver, error)
	GetAllDriversByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	SearchDrivers(searchText, offSet, limit string, companyId int) ([]models.Driver, error)
	GetAllDriversForReport(companyId int) ([]models.Driver, error)
}

type driverService struct {
	repo     repository.DriverRepository
	packages *initpck.Packages
}

func NewDriverService(repo repository.DriverRepository, packages *initpck.Packages) DriverService {
	return &driverService{
		repo:     repo,
		packages: packages,
	}
}

func (service *driverService) GetDriverById(queries *helper.OneQueryParams) (models.Driver, error) {

	driver, err := service.repo.GetDriverById(queries)
	if err != nil {
		return models.Driver{}, err
	}

	return driver, nil
}

func (service *driverService) GetAllDrivers(queries *helper.GetAllQueryParams, companyId int) ([]models.Driver, error) {
	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	drivers, err := service.repo.GetAllDrivers(queries, companyId)
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

func (service *driverService) GetAllDriversByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {
	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	drivers, err := service.repo.GetAllDriversByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedDrivers := service.packages.Sort.SortDataByStructField(orderBy, orderDir, drivers)
		return sortedDrivers, nil
	}

	return drivers, nil
}

func (service *driverService) SearchDrivers(searchText, offSet, limit string, companyId int) ([]models.Driver, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	drivers, err := service.repo.SearchDrivers(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return drivers, nil
}

func (service *driverService) GetAllDriversForReport(companyId int) ([]models.Driver, error) {

	return service.repo.GetAllDriversForReport(companyId)

}
