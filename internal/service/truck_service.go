package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=truck_service.go -destination=mocks/mock_truck_service.go

type TruckService interface {
	GetAllTrucks(queries *helper.GetAllQueryParams, companyId int) ([]models.Truck, error)
	GetAllTrucksByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	GetTruckById(queries *helper.OneQueryParams) (models.Truck, error)
	SearchTrucks(searchText, offSet, limit string, companyId int) ([]models.Truck, error)
	GetAllTrucksForReport(companyId int) ([]models.Truck, error)
}

type truckService struct {
	repo    repository.TruckRepository
	initpck *initpck.Packages
}

func NewTruckService(repo repository.TruckRepository, initpck *initpck.Packages) TruckService {
	return &truckService{repo, initpck}
}

func (service *truckService) GetAllTrucks(queries *helper.GetAllQueryParams, companyId int) ([]models.Truck, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	trucks, err := service.repo.GetAllTrucks(queries, companyId)

	if err != nil {
		return nil, err
	}

	return trucks, nil
}

func (service *truckService) GetAllTrucksByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	trucks, err := service.repo.GetAllTrucksByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedTrucks := service.initpck.Sort.SortDataByStructField(orderBy, orderDir, trucks)
		return sortedTrucks, nil
	}

	return trucks, nil
}

func (service *truckService) GetTruckById(queries *helper.OneQueryParams) (models.Truck, error) {

	truck, err := service.repo.GetTruckById(queries)
	if err != nil {
		return models.Truck{}, err
	}
	return truck, nil
}

func (service *truckService) SearchTrucks(searchText, offSet, limit string, companyId int) ([]models.Truck, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	trucks, err := service.repo.SearchTrucks(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return trucks, nil
}

func (service *truckService) GetAllTrucksForReport(companyId int) ([]models.Truck, error) {

	return service.repo.GetAllTrucksForReport(companyId)

}
