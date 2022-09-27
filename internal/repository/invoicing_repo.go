package repository

import (
	"errors"
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"
	"reflect"
	"strconv"
)

type InvoicingRepository interface {
	GetAllInvoices(queries *helper.PaginationParams) ([]models.Invoicing, error)
	GetInvoicingById(invoicingId string) (models.Invoicing, error)
	SearchAndFilterInvoices(queries *helper.InvoicingFilter) ([]models.Invoicing, error)
	GetAllInvoicesForReport() ([]models.Invoicing, error)
}

type invoicingRepository struct {
	elastic *elastic.Client
}

func NewInvoicingRepository(elastic *elastic.Client) InvoicingRepository {
	return &invoicingRepository{
		elastic: elastic,
	}
}

func (repo *invoicingRepository) GetAllInvoices(queries *helper.PaginationParams) ([]models.Invoicing, error) {

	query := elastic.NewBoolQuery().Must()

	if queries.IsDeleted != "" {
		deleted, _ := strconv.ParseBool(queries.IsDeleted)
		query.Must(elastic.NewMatchQuery("is_deleted", deleted))
	}

	from, _ := strconv.Atoi(queries.OffSet)
	size, _ := strconv.Atoi(queries.Limit)

	res, err := repo.elastic.
		Search("invoicing").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllInvoices", "Cant get all invoices").Error("Error - " + err.Error())
		return []models.Invoicing{}, err
	}

	var ttyp models.Invoicing
	invoices := make([]models.Invoicing, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		invoicing := item.(models.Invoicing)
		invoices = append(invoices, invoicing)
	}

	return invoices, nil
}

func (repo *invoicingRepository) GetInvoicingById(id string) (models.Invoicing, error) {

	query := elastic.NewBoolQuery().Must()

	query.Must(elastic.NewMatchQuery("_id", id))

	res, err := repo.elastic.
		Search("invoicing").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetInvoicingById", "Cant get all invoices by id").Error("Error - " + err.Error())
		return models.Invoicing{}, err
	}

	var ttyp models.Invoicing

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		invoicing := item.(models.Invoicing)
		return invoicing, nil
	}

	return models.Invoicing{}, errors.New("invoicing not found")
}

func (repo *invoicingRepository) SearchAndFilterInvoices(queries *helper.InvoicingFilter) ([]models.Invoicing, error) {

	query := elastic.NewBoolQuery().Must()

	from, _ := strconv.Atoi(queries.OffSet)
	size, _ := strconv.Atoi(queries.Limit)

	if queries.CompanyName != "" {
		query.Must(elastic.NewMatchQuery("company_name", queries.CompanyName))
	}

	if queries.DeliveryFrom != "" {
		deliveryFrom, _ := strconv.ParseInt(queries.DeliveryFrom, 10, 64)
		query.Must(elastic.NewBoolQuery().Filter(elastic.NewRangeQuery("pick_up_date").From(deliveryFrom)))
	}

	if queries.DeliveryTo != "" {
		deliveryTo, _ := strconv.ParseInt(queries.DeliveryTo, 10, 64)
		query.Must(elastic.NewBoolQuery().Filter(elastic.NewRangeQuery("delivery_date").To(deliveryTo)))

	}

	res, err := repo.elastic.
		Search("invoicing").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchAndFilterInvoices", "Cant get all invoices by search and filter").Error("Error - " + err.Error())
		return []models.Invoicing{}, err
	}

	var ttyp models.Invoicing
	invoices := make([]models.Invoicing, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		invoicing := item.(models.Invoicing)
		invoices = append(invoices, invoicing)
	}

	return invoices, nil

}

func (repo *invoicingRepository) GetAllInvoicesForReport() ([]models.Invoicing, error) {

	res, err := repo.elastic.
		Search("invoicing").
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllInvoicesForReport", "Cant get all invoices for report").Error("Error - " + err.Error())
		return []models.Invoicing{}, err
	}

	var ttyp models.Invoicing
	invoices := make([]models.Invoicing, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		invoice := item.(models.Invoicing)
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
