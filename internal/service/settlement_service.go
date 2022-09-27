package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=settlement_service.go -destination=mocks/mock_settlement_service.go

type SettlementService interface {
	GetSettlementById(queries *helper.OneQueryParams) (models.Settlement, error)
	GetAllSettlements(queries *helper.GetAllQueryParams) ([]models.Settlement, error)
	SearchSettlements(searchText, offSet, limit string) ([]models.Settlement, error)
	GetAllSettlementsForReport() ([]models.Settlement, error)
}

type settlementService struct {
	repo     repository.SettlementRepository
	packages *initpck.Packages
}

func NewSettlementService(repo repository.SettlementRepository, packages *initpck.Packages) SettlementService {
	return &settlementService{
		repo:     repo,
		packages: packages,
	}
}

func (service *settlementService) GetAllSettlements(queries *helper.GetAllQueryParams) ([]models.Settlement, error) {
	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	users, err := service.repo.GetAllSettlement(queries)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *settlementService) GetSettlementById(queries *helper.OneQueryParams) (models.Settlement, error) {

	settlement, err := service.repo.GetSettlementById(queries)
	if err != nil {
		return models.Settlement{}, err
	}

	return settlement, nil
}
func (service *settlementService) SearchSettlements(searchText, offSet, limit string) ([]models.Settlement, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	settlements, err := service.repo.SearchSettlements(searchText, offSet, limit)

	if err != nil {
		return nil, err
	}

	return settlements, nil

}

func (service *settlementService) GetAllSettlementsForReport() ([]models.Settlement, error) {

	return service.repo.GetAllSettlementsForReport()

}
