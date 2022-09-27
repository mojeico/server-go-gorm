package service

import (
	"errors"
	"fmt"
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

//go:generate mockgen -source=user_service.go -destination=mocks/mock_user_service.go

//var salt = os.Getenv("salt")

type UserService interface {
	GetAllUsers(queries *helper.GetAllQueryParams, companyId int) ([]models.User, error)
	GetAllUsersByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	GetUserById(queries *helper.OneQueryParams) (models.User, error)
	GetUserByIdForMiddleware(id string) (models.User, error)
	SearchUsers(searchText, offSet, limit string, companyId int) ([]models.User, error)

	SaveUserLogo(*multipart.FileHeader, *gin.Context) (string, error)
}

type userService struct {
	repo     repository.UserRepository
	packages *initpck.Packages
}

func NewUserService(repo repository.UserRepository, packages *initpck.Packages) UserService {
	return &userService{
		repo:     repo,
		packages: packages,
	}
}

func (service *userService) GetAllUsers(queries *helper.GetAllQueryParams, companyId int) ([]models.User, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	users, err := service.repo.GetAllUsers(queries, companyId)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *userService) GetAllUsersByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	users, err := service.repo.GetAllUsersByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedUsers := service.packages.Sort.SortDataByStructField(orderBy, orderDir, users)
		return sortedUsers, nil
	}

	return users, nil
}

func (service *userService) GetUserById(queries *helper.OneQueryParams) (models.User, error) {

	user, err := service.repo.GetUserById(queries)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (service *userService) GetUserByIdForMiddleware(userId string) (models.User, error) {

	user, err := service.repo.GetUserByIdForMiddleware(userId)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (service *userService) SearchUsers(searchText, offSet, limit string, companyId int) ([]models.User, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	users, err := service.repo.SearchUsers(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return users, nil

}

var userFileFormat = []string{".jpg", ".jpeg", ".png", ".svg"}

func (service *userService) SaveUserLogo(fileHeader *multipart.FileHeader, context *gin.Context) (string, error) {

	extension := filepath.Ext(fileHeader.Filename)

	for _, format := range userFileFormat {
		if format == extension {
			newFileName := uuid.New().String() + extension

			err := context.SaveUploadedFile(fileHeader, "upload/userImg/"+newFileName)

			if err != nil {
				logger.ErrorLogger("SaveUserLogo", "Can't save file ").Error("Can't save file ")

				return "", errors.New("can't save file")
			}
			return newFileName, nil
		}
	}

	logger.ErrorLogger("SaveUserLogo", "file extension error").Error(fmt.Sprintf("file extension error - %s", extension))

	return "", errors.New("file extension error")

}
