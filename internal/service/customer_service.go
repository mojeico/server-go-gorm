package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=customer_service.go -destination=mocks/mock_customer_service.go

type CustomerService interface {
	GetCustomerById(queries *helper.OneQueryParams) (models.Customer, error)
	GetAllCustomers(queries *helper.GetAllQueryParams) ([]models.Customer, error)
	SearchCustomers(searchText, offSet, limit string) ([]models.Customer, error)
	GetAllCustomersForReport() ([]models.Customer, error)
}

type customerService struct {
	repo     repository.CustomerRepository
	packages *initpck.Packages
}

func NewCustomerService(repo repository.CustomerRepository, packages *initpck.Packages) CustomerService {
	return &customerService{repo: repo, packages: packages}
}

func (service *customerService) GetCustomerById(queries *helper.OneQueryParams) (models.Customer, error) {

	settlement, err := service.repo.GetCustomerById(queries)
	if err != nil {
		return models.Customer{}, err
	}

	return settlement, nil
}

func (service *customerService) GetAllCustomers(queries *helper.GetAllQueryParams) ([]models.Customer, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	users, err := service.repo.GetAllCustomers(queries)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *customerService) SearchCustomers(searchText, offSet, limit string) ([]models.Customer, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	customers, err := service.repo.SearchCustomers(searchText, offSet, limit)

	if err != nil {
		return nil, err
	}

	return customers, nil

}

func (service *customerService) GetAllCustomersForReport() ([]models.Customer, error) {

	return service.repo.GetAllCustomersForReport()

}
