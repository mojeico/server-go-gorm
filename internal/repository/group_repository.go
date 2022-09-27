package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/trucktrace/pkg/logger"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/helper"
)

type GroupRepository interface {
	GetAllGroups(queries *helper.GetAllQueryParams, companyId int) ([]models.Groups, error)
	GetAllGroupsByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Groups, error)
	GetGroupById(queries *helper.OneQueryParams) (models.Groups, error)
	GetGroupByIdForMiddleware(id string) (models.Groups, error)
	SearchGroups(searchText, offSet, limit string, companyId int) ([]models.Groups, error)
	GetAllGroupsForReport(int) ([]models.Groups, error)
}

type groupRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewGroupRepository(pg *gorm.DB, elastic *elastic.Client) GroupRepository {
	return &groupRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *groupRepository) GetAllGroups(queries *helper.GetAllQueryParams, companyId int) ([]models.Groups, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if queries.QueryField != "" && queries.FieldValue != "" {
		query.Must(elastic.NewMatchQuery(fmt.Sprint(queries.QueryField), queries.FieldValue))
	}

	if queries.IsDeleted != "" {
		deleted, _ := strconv.ParseBool(queries.IsDeleted)
		query.Must(elastic.NewMatchQuery("is_deleted", deleted))
	}

	if queries.IsActive != "" {
		active, _ := strconv.ParseBool(queries.IsActive)
		query.Must(elastic.NewMatchQuery("is_active", active))
	}

	if queries.Status != "" {
		query.Must(
			elastic.NewMultiMatchQuery(queries.Status, "status").Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(queries.OffSet)
	size, _ := strconv.Atoi(queries.Limit)

	res, err := repo.elastic.
		Search("group").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllGroups", "Cant get all groups").Error("Error - " + err.Error())
		return []models.Groups{}, err
	}

	var ttyp models.Groups
	groups := make([]models.Groups, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		group := item.(models.Groups)
		groups = append(groups, group)
	}

	return groups, nil
}

func (repo *groupRepository) GetGroupById(queries *helper.OneQueryParams) (models.Groups, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("_id", queries.Id))

	if queries.IsDeleted != "" {
		deleted, _ := strconv.ParseBool(queries.IsDeleted)
		query.Must(elastic.NewMatchQuery("is_deleted", deleted))
	}

	if queries.IsActive != "" {
		active, _ := strconv.ParseBool(queries.IsActive)
		query.Must(elastic.NewMatchQuery("is_active", active))
	}

	if queries.Status != "" {
		query.Must(
			elastic.NewMultiMatchQuery(queries.Status, "status").Type("phrase_prefix"))
	}

	res, err := repo.elastic.
		Search("group").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetGroupById", "Cant get all group by id").Error("Error - " + err.Error())
		return models.Groups{}, err
	}

	var ttyp models.Groups

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		group := item.(models.Groups)
		return group, nil
	}

	return models.Groups{}, errors.New("group not found")
}

func (repo *groupRepository) GetGroupByIdForMiddleware(groupId string) (models.Groups, error) {

	groupResponse, err := repo.elastic.
		Get().
		Index("group").
		Id(groupId).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetGroupByIdForMiddleware", "Cant get all groups by id for middleware").Error("Error - " + err.Error())
		return models.Groups{}, err
	}

	if groupResponse.Found {
		var group models.Groups

		if err := json.Unmarshal(groupResponse.Source, &group); err != nil {
			logger.ErrorLogger("GetGroupByIdForMiddleware", "Cant unmarshal group").Error("Error - " + err.Error())
			return models.Groups{}, err
		}
		return group, nil
	}
	return models.Groups{}, errors.New("groups not found")

}

func (repo *groupRepository) GetAllGroupsByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Groups, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if queries.QueryField != "" && queries.FieldValue != "" {
		query.Must(elastic.NewMatchQuery(fmt.Sprint(queries.QueryField), queries.FieldValue))
	}

	if queries.IsDeleted != "" {
		deleted, _ := strconv.ParseBool(queries.IsDeleted)
		query.Must(elastic.NewMatchQuery("is_deleted", deleted))
	}

	if queries.IsActive != "" {
		active, _ := strconv.ParseBool(queries.IsActive)
		query.Must(elastic.NewMatchQuery("is_active", active))
	}

	if queries.Status != "" {
		query.Must(
			elastic.NewMultiMatchQuery(queries.Status, "status").Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(queries.OffSet)
	size, _ := strconv.Atoi(queries.Limit)

	res, err := repo.elastic.
		Search("group").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllGroupsByCompanyId", "Cant get all groups by company id").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Groups
	groups := make([]models.Groups, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		group := item.(models.Groups)
		groups = append(groups, group)
	}

	return groups, nil
}

func (repo *groupRepository) SearchGroups(searchText, offSet, limit string, companyId int) ([]models.Groups, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"name",
				"company_id",
				"status",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("group").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchGroups", "Cant get all groups by search").Error("Error - " + err.Error())
		return []models.Groups{}, err
	}

	var ttyp models.Groups
	groups := make([]models.Groups, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		group := item.(models.Groups)
		groups = append(groups, group)
	}

	return groups, nil
}

func (repo *groupRepository) GetAllGroupsForReport(companyId int) ([]models.Groups, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	res, err := repo.elastic.
		Search("group").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllGroupsForReport", "Cant get all groups for report").Error("Error - " + err.Error())
		return []models.Groups{}, err
	}

	var ttyp models.Groups
	groups := make([]models.Groups, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		group := item.(models.Groups)
		groups = append(groups, group)
	}

	return groups, nil
}
