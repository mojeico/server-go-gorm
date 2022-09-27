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

type TruckRepository interface {
	GetAllTrucks(queries *helper.GetAllQueryParams, companyId int) ([]models.Truck, error)
	GetAllTrucksByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Truck, error)
	GetTruckById(queries *helper.OneQueryParams) (models.Truck, error)
	SearchTrucks(searchText, offSet, limit string, companyId int) ([]models.Truck, error)
	GetAllTrucksForReport(int) ([]models.Truck, error)
}

type truckRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewTruckRepository(pg *gorm.DB, elastic *elastic.Client) TruckRepository {
	return &truckRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *truckRepository) GetAllTrucks(queries *helper.GetAllQueryParams, companyId int) ([]models.Truck, error) {

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
		Search("truck").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllTrucks", "Cant get all trucks").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Truck
	trucks := make([]models.Truck, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		truck := item.(models.Truck)
		trucks = append(trucks, truck)
	}

	return trucks, nil
}

func (repo *truckRepository) GetAllTrucksByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Truck, error) {

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
		Search("truck").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllTrucksByCompanyId", "Cant get all trucks by company id").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Truck
	trucks := make([]models.Truck, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		truck := item.(models.Truck)
		trucks = append(trucks, truck)
	}

	return trucks, nil
}

func (repo *truckRepository) GetTruckById(queries *helper.OneQueryParams) (models.Truck, error) {

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
		Search("truck").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetTruckById", "Cant get trucks by id").Error("Error - " + err.Error())
		return models.Truck{}, err
	}

	var ttyp models.Truck

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		truck := item.(models.Truck)
		return truck, nil
	}

	return models.Truck{}, errors.New("truck not found")
}

func (repo *truckRepository) SearchTrucks(searchText, offSet, limit string, companyId int) ([]models.Truck, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"make",
				"model",
				"state",
				"transporter",
				"fuel_card",
				"fuel_type",
				"driver_name",
				"co_driver_name",
				"trailer_name",
				"location",
				"status",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("truck").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchTrucks", "Cant get all trucks by search").Error("Error - " + err.Error())
		return []models.Truck{}, err
	}

	var ttyp models.Truck
	trucks := make([]models.Truck, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		truck := item.(models.Truck)
		trucks = append(trucks, truck)
	}

	return trucks, nil
}

func (repo *truckRepository) GetAllTrucksForReport(companyId int) ([]models.Truck, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	res, err := repo.elastic.
		Search("truck").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllTrucksForReport", "Cant get all trucks for report").Error("Error - " + err.Error())
		return []models.Truck{}, err
	}

	var ttyp models.Truck
	trucks := make([]models.Truck, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		truck := item.(models.Truck)
		trucks = append(trucks, truck)
	}

	return trucks, nil
}
