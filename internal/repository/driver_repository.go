package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/trucktrace/pkg/logger"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/helper"
)

type DriverRepository interface {
	GetDriverById(queries *helper.OneQueryParams) (models.Driver, error)
	GetAllDrivers(queries *helper.GetAllQueryParams, companyId int) ([]models.Driver, error)
	GetAllDriversByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Driver, error)
	SearchDrivers(searchText, offSet, limit string, companyId int) ([]models.Driver, error)
	GetAllDriversForReport(companyId int) ([]models.Driver, error)
}

type driverRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewDriverRepository(pg *gorm.DB, elastic *elastic.Client) DriverRepository {
	return &driverRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *driverRepository) GetAllDrivers(queries *helper.GetAllQueryParams, companyId int) ([]models.Driver, error) {

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
		Search("driver").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllDrivers", "Cant get all drivers").Error("Error - " + err.Error())
		return nil, err
	}

	drivers := make([]models.Driver, 0)

	for i := 0; i < int(res.TotalHits()); i++ {
		var driver models.Driver
		driverByte := res.Hits.Hits[i].Source
		_ = json.Unmarshal(driverByte, &driver)
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

func (repo *driverRepository) GetDriverById(queries *helper.OneQueryParams) (models.Driver, error) {

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
		Search("driver").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetDriverById", "Cant get driver by id").Error("Error - " + err.Error())
		return models.Driver{}, err
	}

	if res.TotalHits() > 0 {
		var driver models.Driver
		_ = json.Unmarshal(res.Hits.Hits[0].Source, &driver)

		return driver, nil
	}

	return models.Driver{}, errors.New("driver not found")
}

func (repo *driverRepository) GetAllDriversByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Driver, error) {

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
		Search("driver").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllDriversByCompanyId", "Cant get all drivers by company id ").Error("Error - " + err.Error())
		return nil, err
	}

	drivers := make([]models.Driver, 0)

	for i := 0; i < int(res.TotalHits()); i++ {
		var driver models.Driver

		driverByte := res.Hits.Hits[i].Source

		_ = json.Unmarshal(driverByte, &driver)
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

func (repo *driverRepository) SearchDrivers(searchText, offSet, limit string, companyId int) ([]models.Driver, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"first_name",
				"last_name",
				"address",
				"email",
				"phone",
				"city",
				"state",
				"country",
				"phone",
				"gender",
				"pay_to",
				"pay_method",
				"pay_to_count",
				"pay_to_owner",
				"fuel_cards",
				"transponder",
				"ssn",
				"cdl_classification",
				"cdl_endorsements",
				"tax_form",
				"restrictions",
				"emergency",
				"emergency_phone",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("driver").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.InfoLogger("SearchDrivers").Error("Error - " + err.Error())
		return []models.Driver{}, err
	}

	drivers := make([]models.Driver, 0)

	for i := 0; i < int(res.TotalHits()); i++ {
		var driver models.Driver

		driverByte := res.Hits.Hits[i].Source

		_ = json.Unmarshal(driverByte, &driver)
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

func (repo *driverRepository) GetAllDriversForReport(companyId int) ([]models.Driver, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	res, err := repo.elastic.
		Search("driver").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllDriversForReport", "Cant get all drivers for report").Error("Error - " + err.Error())
		return []models.Driver{}, err
	}

	drivers := make([]models.Driver, 0)

	for i := 0; i < int(res.TotalHits()); i++ {
		var driver models.Driver

		driverByte := res.Hits.Hits[i].Source
		_ = json.Unmarshal(driverByte, &driver)

		drivers = append(drivers, driver)
	}

	return drivers, nil
}
