package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

type InvoicingService interface {
	GetAllInvoices(query *helper.PaginationParams) ([]models.Invoicing, error)
	GetInvoicingById(invoicingId string) (models.Invoicing, error)
	SearchAndFilterInvoices(query *helper.InvoicingFilter) ([]models.Invoicing, error)
	GetAllInvoicesForReport() ([]models.Invoicing, error)
}

type invoicingService struct {
	repo     repository.InvoicingRepository
	packages *initpck.Packages
}

func NewInvoicingService(repo repository.InvoicingRepository, packages *initpck.Packages) InvoicingService {
	return &invoicingService{
		repo:     repo,
		packages: packages,
	}
}

func (service *invoicingService) GetAllInvoices(query *helper.PaginationParams) ([]models.Invoicing, error) {

	query.OffSet, query.Limit = helper.CheckLimit(query.OffSet, query.Limit)

	invoices, err := service.repo.GetAllInvoices(query)

	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (service *invoicingService) GetInvoicingById(invoicingId string) (models.Invoicing, error) {

	invoicing, err := service.repo.GetInvoicingById(invoicingId)
	if err != nil {
		return models.Invoicing{}, err
	}

	return invoicing, nil
}

func (service *invoicingService) SearchAndFilterInvoices(query *helper.InvoicingFilter) ([]models.Invoicing, error) {

	query.OffSet, query.Limit = helper.CheckLimit(query.OffSet, query.Limit)

	invoices, err := service.repo.SearchAndFilterInvoices(query)

	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (service *invoicingService) GetAllInvoicesForReport() ([]models.Invoicing, error) {

	return service.repo.GetAllInvoicesForReport()

}
