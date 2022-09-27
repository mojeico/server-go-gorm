package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=order_service.go -destination=mocks/mock_order_service.go

type OrderService interface {
	GetAllOrders(queries *helper.GetAllQueryParams, companyId int) ([]models.Order, error)
	GetAllOrdersByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	GetOrderById(queries *helper.OneQueryParams) (models.Order, error)
	GetOrderByDriverId(queries *helper.OneQueryParams) ([]models.Order, error)
	GetOrderByTruckId(queries *helper.OneQueryParams) ([]models.Order, error)
	GetOrderByTrailerId(queries *helper.OneQueryParams) ([]models.Order, error)
	SearchOrders(searchText, offSet, limit string, companyId int) ([]models.Order, error)
	GetAllOrdersForReport(companyId int) ([]models.Order, error)
	GetExtraPaysByOrderId(queries *helper.OneQueryParams) ([]models.ExtraPay, error)
}

type orderService struct {
	repo    repository.OrderRepository
	initpck *initpck.Packages
}

func NewOrderService(repo repository.OrderRepository, initpck *initpck.Packages) OrderService {
	return &orderService{repo, initpck}
}

func (service *orderService) GetAllOrders(queries *helper.GetAllQueryParams, companyId int) ([]models.Order, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	orders, err := service.repo.GetAllOrders(queries, companyId)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) GetAllOrdersByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	orders, err := service.repo.GetAllOrdersByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedOrders := service.initpck.Sort.SortDataByStructField(orderBy, orderDir, orders)
		return sortedOrders, nil
	}

	return orders, nil
}

func (service *orderService) GetOrderById(queries *helper.OneQueryParams) (models.Order, error) {

	order, err := service.repo.GetOrderById(queries)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (service *orderService) GetOrderByDriverId(queries *helper.OneQueryParams) ([]models.Order, error) {
	order, err := service.repo.GetOrdersByDriverId(queries)
	if err != nil {
		return []models.Order{}, err
	}
	return order, nil
}

func (service *orderService) GetOrderByTruckId(queries *helper.OneQueryParams) ([]models.Order, error) {
	order, err := service.repo.GetOrdersByTruckId(queries)
	if err != nil {
		return []models.Order{}, err
	}
	return order, nil
}

func (service *orderService) GetOrderByTrailerId(queries *helper.OneQueryParams) ([]models.Order, error) {
	order, err := service.repo.GetOrdersByTrailerId(queries)
	if err != nil {
		return []models.Order{}, err
	}
	return order, nil
}

func (service *orderService) SearchOrders(searchText, offSet, limit string, companyId int) ([]models.Order, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	orders, err := service.repo.SearchOrders(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *orderService) GetAllOrdersForReport(companyId int) ([]models.Order, error) {
	return service.repo.GetAllOrdersForReport(companyId)
}

func (service *orderService) GetExtraPaysByOrderId(queries *helper.OneQueryParams) ([]models.ExtraPay, error) {
	return service.repo.GetExtraPaysByOrderId(queries)
}
