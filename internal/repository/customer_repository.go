package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/trucktrace/pkg/logger"
	"reflect"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/pkg/helper"

	"github.com/jinzhu/gorm"
	"github.com/trucktrace/internal/models"
)

type CustomerRepository interface {
	GetCustomerById(queries *helper.OneQueryParams) (models.Customer, error)
	GetAllCustomers(queries *helper.GetAllQueryParams) ([]models.Customer, error)
	SearchCustomers(searchText, offSet, limit string) ([]models.Customer, error)
	GetAllCustomersForReport() ([]models.Customer, error)
}

type customerRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewCustomerRepository(pg *gorm.DB, elastic *elastic.Client) CustomerRepository {
	return &customerRepository{
		pg:      pg,
		elastic: elastic,
	}
}

var ctx = context.Background()

func (repo customerRepository) GetCustomerById(queries *helper.OneQueryParams) (models.Customer, error) {

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
		Search("customer").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetCustomerById", "Cant get customer by id").Error("Error - " + err.Error())
		return models.Customer{}, err
	}

	var ttyp models.Customer

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		customer := item.(models.Customer)
		return customer, nil
	}

	return models.Customer{}, errors.New("customer not found")

}

func (repo *customerRepository) GetAllCustomers(queries *helper.GetAllQueryParams) ([]models.Customer, error) {

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
		Search("customer").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllCustomers", "Cant get all customers").Error("Error- " + err.Error())
		return []models.Customer{}, err
	}

	var ttyp models.Customer
	customers := make([]models.Customer, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		customer := item.(models.Customer)
		customers = append(customers, customer)
	}

	return customers, nil

}

func (repo *customerRepository) SearchCustomers(searchText, offSet, limit string) ([]models.Customer, error) {

	query := elastic.NewBoolQuery().Must()

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"address",
				"city",
				"country",
				"phone_number",
				"fax_number",
				"legal_name",
				"dot_number",
				"billing_address",
				"billing_method",
				"billing_type",
				"billing_email",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("customer").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchCustomers", "Cant get all customers by search").Error("Error - " + err.Error())
		return []models.Customer{}, err
	}

	var ttyp models.Customer
	customers := make([]models.Customer, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		customer := item.(models.Customer)
		customers = append(customers, customer)
	}

	return customers, nil
}

func (repo *customerRepository) GetAllCustomersForReport() ([]models.Customer, error) {

	res, err := repo.elastic.
		Search("customer").
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllCustomersForReport", "Cant get all customers for report").Error("Error - " + err.Error())
		return []models.Customer{}, err
	}

	var ttyp models.Customer
	customers := make([]models.Customer, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		customer := item.(models.Customer)
		customers = append(customers, customer)
	}

	return customers, nil
}
