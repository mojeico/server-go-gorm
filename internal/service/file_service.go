package service

import (
	"errors"
	"fmt"
	"github.com/trucktrace/pkg/logger"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

type FileService interface {
	UploadFile(filename string) (string, string, error)
	GetFileById(queries *helper.OneQueryParams) (models.File, error)
	GetLastFileByOwnerId(ownerType string, queries *helper.OneQueryParams) (models.File, error)
	GetAllFilesByOwnerId(ownerType string, queries *helper.GetAllQueryParams) ([]models.File, error)
	GetAllFiles(queries *helper.GetAllQueryParams) ([]models.File, error)
	SearchFiles(searchText, offSet, limit string) ([]models.File, error)
}

type fileService struct {
	repo     repository.FileRepository
	packages *initpck.Packages
}

func NewFileService(repo repository.FileRepository, packages *initpck.Packages) FileService {
	return &fileService{repo, packages}
}

var fileFormat = []string{".pdf", ".xlsx", ".xlsm", ".xls", ".csv", ".jpg", ".jpeg", ".png"}

func (service *fileService) UploadFile(filename string) (string, string, error) {

	extension := filepath.Ext(filename)

	for _, format := range fileFormat {
		if format == extension {
			newFileName := uuid.New().String() + extension
			return newFileName, extension, nil
		}
	}

	logger.ErrorLogger("UploadFile", "file extension error").Error(fmt.Sprintf("file extension error - %s", extension))

	return "", "", errors.New("file extension error")
}

func (service *fileService) GetFileById(queries *helper.OneQueryParams) (models.File, error) {
	file, err := service.repo.GetFileById(queries)

	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

func (service *fileService) GetLastFileByOwnerId(ownerType string, queries *helper.OneQueryParams) (models.File, error) {

	if !checkOwnerType(ownerType) {
		logger.ErrorLogger("GetLastFileByOwnerId", "wrong owner type").Error(fmt.Sprintf("wrong owner type - %s", ownerType))
		return models.File{}, errors.New("wrong owner type")
	}

	file, err := service.repo.GetLastFileByOwnerId(ownerType, queries)

	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

func (service *fileService) GetAllFilesByOwnerId(ownerType string, queries *helper.GetAllQueryParams) ([]models.File, error) {

	if !checkOwnerType(ownerType) {
		logger.ErrorLogger("GetAllFilesByOwnerId", "wrong owner type").Error(fmt.Sprintf("wrong owner type - %s", ownerType))
		return nil, errors.New("wrong owner type")
	}

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	files, err := service.repo.GetAllFilesByOwnerId(ownerType, queries)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (service *fileService) GetAllFiles(queries *helper.GetAllQueryParams) ([]models.File, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	files, err := service.repo.GetAllFiles(queries)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (service *fileService) SearchFiles(searchText, offSet, limit string) ([]models.File, error) {
	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	safeties, err := service.repo.SearchFiles(searchText, offSet, limit)

	if err != nil {
		return nil, err
	}

	return safeties, nil
}

func checkOwnerType(owner string) bool {
	var ownerTypes = []string{"driver", "order", "safety", "customer", "trailer", "truck"}

	for _, ot := range ownerTypes {
		if owner == ot {
			return true
		}
	}

	return false
}
