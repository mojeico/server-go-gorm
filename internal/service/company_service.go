package service

import (
	"errors"
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"

	initpck "github.com/trucktrace/pkg/main"
)

type CompanyService interface {
	GetCompanyById(queries *helper.OneQueryParams) (models.Company, error)
	GetAllCompanies(queries *helper.GetAllQueryParams) ([]models.Company, error)
	SearchCompany(searchText, offSet, limit string) ([]models.Company, error)
	GetAllCompaniesForReport() ([]models.Company, error)

	SaveOrderLogo(*multipart.FileHeader, *gin.Context) (string, error)
}

type companyService struct {
	repo     repository.CompanyRepository
	packages *initpck.Packages
}

func NewCompanyService(repo repository.CompanyRepository, packages *initpck.Packages) CompanyService {
	return &companyService{
		repo:     repo,
		packages: packages,
	}
}

func (service companyService) GetCompanyById(queries *helper.OneQueryParams) (models.Company, error) {

	settlement, err := service.repo.GetCompanyById(queries)
	if err != nil {
		return models.Company{}, err
	}

	return settlement, nil
}

func (service companyService) GetAllCompanies(queries *helper.GetAllQueryParams) ([]models.Company, error) {
	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	companies, err := service.repo.GetAllCompanies(queries)

	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (service companyService) SearchCompany(searchText, offSet, limit string) ([]models.Company, error) {
	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	companies, err := service.repo.SearchCompany(searchText, offSet, limit)

	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (service *companyService) GetAllCompaniesForReport() ([]models.Company, error) {

	return service.repo.GetAllCompaniesForReport()
}

var companyFileFormat = []string{".jpg", ".jpeg", ".png", ".svg"}

func (service *companyService) SaveOrderLogo(fileHeader *multipart.FileHeader, context *gin.Context) (string, error) {

	extension := filepath.Ext(fileHeader.Filename)

	for _, format := range companyFileFormat {
		if format == extension {
			newFileName := uuid.New().String() + extension

			err := context.SaveUploadedFile(fileHeader, "upload/companyImg/"+newFileName)

			if err != nil {
				logger.ErrorLogger("SaveOrderLogo", "Can't save file").Error("Error - " + err.Error())
				return "", errors.New("can't save file")
			}
			return newFileName, nil
		}
	}

	return "", errors.New("file extension error")

}
