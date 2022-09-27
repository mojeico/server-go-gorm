package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

type ChargesService interface {
	GetAllCharges(queries *helper.GetAllQueryParams) ([]models.Charges, error)
	GetAllChargesBySettlementId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error)
	GetAllChargesByOrderId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error)
	GetChargeById(queries *helper.OneQueryParams) (models.Charges, error)
	SearchCharge(searchText, offSet, limit string) ([]models.Charges, error)
	GetAllChargesForReport() ([]models.Charges, error)
}

type chargesService struct {
	repo     repository.ChargesRepository
	packages *initpck.Packages
}

func NewChargesService(repo repository.ChargesRepository, packages *initpck.Packages) ChargesService {
	return &chargesService{repo: repo, packages: packages}
}

func (service *chargesService) GetAllCharges(queries *helper.GetAllQueryParams) ([]models.Charges, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	charges, err := service.repo.GetAllCharges(queries)

	if err != nil {
		return nil, err
	}

	return charges, nil
}

func (service *chargesService) GetAllChargesBySettlementId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error) {

	queries.OffSet, queries.OffSet = helper.CheckLimit(queries.OffSet, queries.OffSet)

	charges, err := service.repo.GetAllChargesBySettlementId(queries)

	if err != nil {
		return nil, err
	}

	return charges, nil
}

func (service *chargesService) GetAllChargesByOrderId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error) {

	queries.OffSet, queries.OffSet = helper.CheckLimit(queries.OffSet, queries.OffSet)

	charges, err := service.repo.GetAllChargesByOrderId(queries)

	if err != nil {
		return nil, err
	}

	return charges, nil
}

func (service *chargesService) GetChargeById(queries *helper.OneQueryParams) (models.Charges, error) {

	charges, err := service.repo.GetChargeById(queries)
	if err != nil {
		return models.Charges{}, err
	}

	return charges, nil
}

func (service *chargesService) SearchCharge(searchText, offSet, limit string) ([]models.Charges, error) {

	offSet, limit = helper.CheckLimit(offSet, limit)

	charges, err := service.repo.SearchCharge(searchText, offSet, limit)

	if err != nil {
		return nil, err
	}

	return charges, nil
}

func (service *chargesService) GetAllChargesForReport() ([]models.Charges, error) {

	return service.repo.GetAllChargesForReport()

}
