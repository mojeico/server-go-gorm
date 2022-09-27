package repository

import (
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

type SafetyRepository interface {
	GetSafetyById(companyId string, queries *helper.OneQueryParams) (models.Safety, error)
	GetAllSafetiesByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Safety, error)
	GetAllSafeties(queries *helper.GetAllQueryParams, companyId int) ([]models.Safety, error)
	SearchSafety(searchText, offSet, limit string, id int) ([]models.Safety, error)
	GetAllSafetyByType(safetyType, offSet, limit string, companyId int) ([]models.Safety, error)
	GetAllSafetiesForReport(int) ([]models.Safety, error)
}

type safetyRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewSafetyRepository(pg *gorm.DB, elastic *elastic.Client) SafetyRepository {
	return &safetyRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *safetyRepository) GetAllSafetiesByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Safety, error) {

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
		Search("safety").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllSafetiesByCompanyId", "Cant get all safeties by company id ").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Safety
	safeties := make([]models.Safety, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		safety := item.(models.Safety)
		safeties = append(safeties, safety)
	}

	return safeties, nil
}

func (repo safetyRepository) GetSafetyById(companyId string, queries *helper.OneQueryParams) (models.Safety, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("_id", queries.Id))
	query.Must(elastic.NewMatchQuery("company_id", companyId))

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
		Search("safety").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetSafetyById", "Cant get safety by id").Error("Error - " + err.Error())
		return models.Safety{}, err
	}

	var ttyp models.Safety

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		safety := item.(models.Safety)
		return safety, nil
	}

	return models.Safety{}, errors.New("safety not found")
}

func (repo safetyRepository) GetAllSafeties(queries *helper.GetAllQueryParams, companyId int) ([]models.Safety, error) {

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
		Search("safety").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllSafeties", "Cant get all safeties").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Safety
	safeties := make([]models.Safety, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		safety := item.(models.Safety)
		safeties = append(safeties, safety)
	}

	return safeties, nil
}

func (repo *safetyRepository) SearchSafety(searchText, offSet, limit string, companyId int) ([]models.Safety, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"file_name",
				"safety_type",
				"comments",
				"status",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("safety").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchSafety", "Cant get all safeties by search").Error("Error - " + err.Error())
		return []models.Safety{}, err
	}

	var ttyp models.Safety
	safeties := make([]models.Safety, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		safety := item.(models.Safety)
		safeties = append(safeties, safety)
	}

	return safeties, nil
}

func (repo *safetyRepository) GetAllSafetyByType(safetyType, offSet, limit string, companyId int) ([]models.Safety, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("safety_type", safetyType))
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("safety").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllSafetyByType", "Cant get all safeties by type").Error("Error - " + err.Error())
		return []models.Safety{}, err
	}

	var ttyp models.Safety
	safeties := make([]models.Safety, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		safety := item.(models.Safety)
		safeties = append(safeties, safety)
	}

	return safeties, nil
}

func (repo *safetyRepository) GetAllSafetiesForReport(companyId int) ([]models.Safety, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	res, err := repo.elastic.
		Search("safety").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllSafetiesForReport", "Cant get all safeties for report").Error("Error - " + err.Error())
		return []models.Safety{}, err
	}

	var ttyp models.Safety
	safeties := make([]models.Safety, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		safety := item.(models.Safety)
		safeties = append(safeties, safety)
	}

	return safeties, nil
}
