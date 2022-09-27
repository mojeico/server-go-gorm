package repository

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/trucktrace/pkg/logger"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/helper"
)

type OrderRepository interface {
	GetAllOrders(queries *helper.GetAllQueryParams, companyId int) ([]models.Order, error)
	GetAllOrdersByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Order, error)
	GetOrderById(queries *helper.OneQueryParams) (models.Order, error)
	GetOrdersByDriverId(queries *helper.OneQueryParams) ([]models.Order, error)
	GetOrdersByTruckId(queries *helper.OneQueryParams) ([]models.Order, error)
	GetOrdersByTrailerId(queries *helper.OneQueryParams) ([]models.Order, error)
	SearchOrders(searchText, offSet, limit string, companyId int) ([]models.Order, error)
	GetAllOrdersForReport(companyId int) ([]models.Order, error)
	GetExtraPaysByOrderId(queries *helper.OneQueryParams) ([]models.ExtraPay, error)
}

type orderRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewOrderRepository(pg *gorm.DB, elastic *elastic.Client) OrderRepository {
	return &orderRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *orderRepository) GetAllOrders(queries *helper.GetAllQueryParams, companyId int) ([]models.Order, error) {

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
		Search("order").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllOrders", "Cant get all orders").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil

}

func (repo *orderRepository) GetAllOrdersByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.Order, error) {

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
		Search("order").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllOrdersByCompanyId", "Cant get all orders by company id").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) GetOrderById(queries *helper.OneQueryParams) (models.Order, error) {

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
		Search("order").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetOrderById", "Cant get  order by id").Error("Error - " + err.Error())
		return models.Order{}, err
	}

	var ttyp models.Order

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		return order, nil
	}

	return models.Order{}, nil
}

func (repo *orderRepository) GetOrdersByDriverId(queries *helper.OneQueryParams) ([]models.Order, error) {
	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("driver_id", queries.Id))

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
		Search("order").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetOrdersByDriverId", "Cant get all orders by driver id").Error("Error - " + err.Error())
		return []models.Order{}, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) GetOrdersByTruckId(queries *helper.OneQueryParams) ([]models.Order, error) {
	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("truck_id", queries.Id))

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
		Search("order").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetOrdersByTruckId", "Cant get all orders by truck id").Error("Error - " + err.Error())
		return []models.Order{}, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) GetOrdersByTrailerId(queries *helper.OneQueryParams) ([]models.Order, error) {
	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("trailer_id", queries.Id))

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
		Search("order").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetOrdersByTrailerId", "Cant get all orders by trailer id ").Error("Error - " + err.Error())
		return []models.Order{}, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) SearchOrders(searchText, offSet, limit string, companyId int) ([]models.Order, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"load_number",
				"pickup_number",
				"delivery_number",
				"shipper",
				"shipper_from_location",
				"shipper_phone",
				"consignee",
				"consignee_to_location",
				"consignee_phone",
				"seal_number",
				"bol_number",
				"commodity",
				"equipment_type",
				"pieces",
				"invoicing_company",
				"billing_method",
				"billing_type",
				"driver_name",
				"external_notes",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("order").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchOrders", "Cant get all orders by search").Error("Error - " + err.Error())
		return []models.Order{}, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) GetAllOrdersForReport(companyId int) ([]models.Order, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	res, err := repo.elastic.
		Search("order").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllOrdersForReport", "Cant get all orders for report").Error("Error - " + err.Error())
		return []models.Order{}, err
	}

	var ttyp models.Order
	orders := make([]models.Order, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		order := item.(models.Order)
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *orderRepository) GetExtraPaysByOrderId(queries *helper.OneQueryParams) ([]models.ExtraPay, error) {

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
		query.Must(
			elastic.NewMultiMatchQuery(queries.Status, "status").Type("phrase_prefix"))
	}

	res, err := repo.elastic.
		Search("extrapay").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetExtraPaysByOrderId", "Cant get all extra pay by order id ").Error("Error - " + err.Error())
		return []models.ExtraPay{}, err
	}

	var ttyp models.ExtraPay
	extraPays := make([]models.ExtraPay, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		extraPay := item.(models.ExtraPay)
		extraPays = append(extraPays, extraPay)
	}

	return extraPays, nil
}
