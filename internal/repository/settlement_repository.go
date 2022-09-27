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

type SettlementRepository interface {
	GetSettlementById(queries *helper.OneQueryParams) (models.Settlement, error)
	GetAllSettlement(queries *helper.GetAllQueryParams) ([]models.Settlement, error)
	SearchSettlements(searchText, offSet, limit string) ([]models.Settlement, error)
	GetAllSettlementsForReport() ([]models.Settlement, error)
}

type settlementRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewSettlementRepository(pg *gorm.DB, elastic *elastic.Client) SettlementRepository {
	return &settlementRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *settlementRepository) GetSettlementById(queries *helper.OneQueryParams) (models.Settlement, error) {

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
		Search("settlement").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetSettlementById", "Cant get settlement by id ").Error("Error - " + err.Error())
		return models.Settlement{}, err
	}

	var ttyp models.Settlement

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		settlement := item.(models.Settlement)
		return settlement, nil
	}

	return models.Settlement{}, errors.New("settlement not found")
}

func (repo *settlementRepository) GetAllSettlement(queries *helper.GetAllQueryParams) ([]models.Settlement, error) {

	query := elastic.NewBoolQuery().Must()

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
		Search().
		Index("settlement").
		From(from).
		Size(size).
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllSettlement", "Cant get all settlements").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Settlement
	settlements := make([]models.Settlement, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		settlement := item.(models.Settlement)
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

func (repo *settlementRepository) SearchSettlements(searchText, offSet, limit string) ([]models.Settlement, error) {

	query := elastic.NewBoolQuery().Must()

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"invoicing_company",
				"driver_name",
				"status",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("settlement").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchSettlements", "Cant get all settlements by search").Error("Error - " + err.Error())

		return []models.Settlement{}, err
	}

	var ttyp models.Settlement
	settlements := make([]models.Settlement, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		settlement := item.(models.Settlement)
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

func (repo *settlementRepository) GetAllSettlementsForReport() ([]models.Settlement, error) {

	res, err := repo.elastic.
		Search("settlement").
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllSettlementsForReport", "Cant get all settlements for report").Error("Error - " + err.Error())
		return []models.Settlement{}, err
	}

	var ttyp models.Settlement
	settlements := make([]models.Settlement, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		settlement := item.(models.Settlement)
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}
