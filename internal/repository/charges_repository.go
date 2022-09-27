package repository

import (
	"errors"
	"fmt"
	"github.com/trucktrace/pkg/logger"
	"reflect"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/helper"
)

type ChargesRepository interface {
	GetAllCharges(queries *helper.GetAllQueryParams) ([]models.Charges, error)
	GetAllChargesBySettlementId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error)
	GetAllChargesByOrderId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error)
	GetChargeById(queries *helper.OneQueryParams) (models.Charges, error)
	SearchCharge(searchText, offSet, limit string) ([]models.Charges, error)

	GetAllChargesForReport() ([]models.Charges, error)
}

type chargesRepository struct {
	elastic *elastic.Client
}

func NewChargesRepository(elastic *elastic.Client) ChargesRepository {
	return &chargesRepository{
		elastic: elastic,
	}
}

func (repo *chargesRepository) GetAllCharges(queries *helper.GetAllQueryParams) ([]models.Charges, error) {

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
		query.Must(elastic.NewMatchQuery("status", queries.Status))
	}

	from, _ := strconv.Atoi(queries.OffSet)
	size, _ := strconv.Atoi(queries.Limit)

	res, err := repo.elastic.
		Search("charges").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.InfoLogger("GetAllCharges").Error("Error - " + err.Error())
		return []models.Charges{}, err
	}

	var ttyp models.Charges
	charges := make([]models.Charges, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		charge := item.(models.Charges)
		charges = append(charges, charge)
	}

	return charges, nil
}

func (repo *chargesRepository) GetAllChargesBySettlementId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error) {

	query := elastic.NewBoolQuery().Must()

	query.Must(elastic.NewMatchQuery("settlement_id", queries.Id))

	if queries.IsDeleted != "" {
		deleted, _ := strconv.ParseBool(queries.IsDeleted)
		query.Must(elastic.NewMatchQuery("is_deleted", deleted))
	}

	if queries.IsActive != "" {
		active, _ := strconv.ParseBool(queries.IsActive)
		query.Must(elastic.NewMatchQuery("is_active", active))
	}

	if queries.Status != "" {
		query.Must(elastic.NewMatchQuery("status", queries.Status))
	}

	res, err := repo.elastic.
		Search("charges").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllChargesBySettlementId", "Cant get all charges by settlement id").Error("Error - " + err.Error())
		return []models.Charges{}, err
	}

	var ttyp models.Charges
	charges := make([]models.Charges, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		charge := item.(models.Charges)
		charges = append(charges, charge)
	}

	return charges, nil

}

func (repo *chargesRepository) GetAllChargesByOrderId(queries *helper.GetAllQueryParamsWithId) ([]models.Charges, error) {

	query := elastic.NewBoolQuery().Must()

	query.Must(elastic.NewMatchQuery("order_id", queries.Id))

	if queries.IsDeleted != "" {
		deleted, _ := strconv.ParseBool(queries.IsDeleted)
		query.Must(elastic.NewMatchQuery("is_deleted", deleted))
	}

	if queries.IsActive != "" {
		active, _ := strconv.ParseBool(queries.IsActive)
		query.Must(elastic.NewMatchQuery("is_active", active))
	}

	if queries.Status != "" {
		query.Must(elastic.NewMatchQuery("status", queries.Status))
	}

	res, err := repo.elastic.
		Search("charges").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllChargesByOrderId", "Cant get all charges by order id").Error("Error - " + err.Error())
		return []models.Charges{}, err
	}

	var ttyp models.Charges
	charges := make([]models.Charges, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		charge := item.(models.Charges)
		charges = append(charges, charge)
	}

	return charges, nil
}

func (repo *chargesRepository) GetChargeById(queries *helper.OneQueryParams) (models.Charges, error) {

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
		query.Must(elastic.NewMatchQuery("status", queries.Status))
	}

	res, err := repo.elastic.
		Search("charges").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetChargeById", "Cant charge by  id").Error("Error - " + err.Error())
		return models.Charges{}, err
	}

	var ttyp models.Charges

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		charge := item.(models.Charges)
		return charge, nil
	}

	return models.Charges{}, errors.New("charge not found")
}

func (repo *chargesRepository) SearchCharge(searchText, offSet, limit string) ([]models.Charges, error) {

	query := elastic.NewBoolQuery().Must()

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"driver_name",
				"company_name",
				"type_deductions",
				"description",
				"status").Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("charges").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchCharge", "Cant get all charges by serarch").Error("Error - " + err.Error())
		return []models.Charges{}, err
	}

	var ttyp models.Charges
	charges := make([]models.Charges, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		charge := item.(models.Charges)
		charges = append(charges, charge)
	}

	return charges, nil
}

func (repo *chargesRepository) GetAllChargesForReport() ([]models.Charges, error) {

	res, err := repo.elastic.
		Search("charges").
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllChargesForReport", "Cant get all charges for report").Error("Error - " + err.Error())
		return []models.Charges{}, err
	}

	var ttyp models.Charges
	charges := make([]models.Charges, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		charge := item.(models.Charges)
		charges = append(charges, charge)
	}

	return charges, nil
}
