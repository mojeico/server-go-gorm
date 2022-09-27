package service

import (
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/helper"
	initpck "github.com/trucktrace/pkg/main"
)

//go:generate mockgen -source=group_service.go -destination=mocks/mock_group_service.go

type GroupService interface {
	GetAllGroupsByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error)
	GetGroupById(queries *helper.OneQueryParams) (models.Groups, error)
	GetGroupByIdForMiddlewares(id string) (models.Groups, error)

	GetAllGroups(queries *helper.GetAllQueryParams, companyId int) ([]models.Groups, error)
	SearchGroups(searchText, offSet, limit string, companyId int) ([]models.Groups, error)
	GetAllGroupsForReport(companyId int) ([]models.Groups, error)
}

type groupService struct {
	repo     repository.GroupRepository
	packages *initpck.Packages
}

func NewGroupService(repo repository.GroupRepository, packages *initpck.Packages) GroupService {
	return &groupService{
		repo:     repo,
		packages: packages,
	}
}

func (service *groupService) GetAllGroupsByCompanyId(companyId string, queries *helper.GetAllQueryParams, orderBy string, orderDir byte) (interface{}, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	groups, err := service.repo.GetAllGroupsByCompanyId(companyId, queries)
	if err != nil {
		return nil, err
	}

	if orderBy != "" {
		sortedGroups := service.packages.Sort.SortDataByStructField(orderBy, orderDir, groups)
		return sortedGroups, nil
	}

	return groups, nil
}

func (service *groupService) GetAllGroups(queries *helper.GetAllQueryParams, companyId int) ([]models.Groups, error) {

	if queries.OffSet == "" {
		queries.OffSet = "0"
	}

	if queries.Limit == "" {
		queries.Limit = "20"
	}

	groups, err := service.repo.GetAllGroups(queries, companyId)

	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (service *groupService) GetGroupById(queries *helper.OneQueryParams) (models.Groups, error) {

	group, err := service.repo.GetGroupById(queries)

	if err != nil {
		return models.Groups{}, err
	}

	return group, nil
}

func (service *groupService) GetGroupByIdForMiddlewares(id string) (models.Groups, error) {

	group, err := service.repo.GetGroupByIdForMiddleware(id)

	if err != nil {
		return models.Groups{}, err
	}

	return group, nil
}

func (service *groupService) SearchGroups(searchText, offSet, limit string, companyId int) ([]models.Groups, error) {

	if offSet == "" {
		offSet = "0"
	}

	if limit == "" {
		limit = "20"
	}

	groups, err := service.repo.SearchGroups(searchText, offSet, limit, companyId)

	if err != nil {
		return nil, err
	}

	return groups, nil

}

func (service *groupService) GetAllGroupsForReport(companyId int) ([]models.Groups, error) {

	return service.repo.GetAllGroupsForReport(companyId)

}
