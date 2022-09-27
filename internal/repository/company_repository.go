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

type CompanyRepository interface {
	GetCompanyById(queries *helper.OneQueryParams) (models.Company, error)
	GetAllCompanies(queries *helper.GetAllQueryParams) ([]models.Company, error)
	SearchCompany(searchText, offSet, limit string) ([]models.Company, error)
	GetAllCompaniesForReport() ([]models.Company, error)
}

type companyRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewCompanyRepository(pg *gorm.DB, elastic *elastic.Client) CompanyRepository {
	return &companyRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *companyRepository) GetCompanyById(queries *helper.OneQueryParams) (models.Company, error) {
	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("_id", queries.Id))

	res, err := repo.elastic.
		Search("company").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetCompanyById", "Cant get all company by id").Error("Error - " + err.Error())
		return models.Company{}, err
	}

	var ttyp models.Company

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		company := item.(models.Company)
		return company, nil
	}

	return models.Company{}, errors.New("company not found")
}

func (repo *companyRepository) GetAllCompanies(queries *helper.GetAllQueryParams) ([]models.Company, error) {
	query := elastic.NewBoolQuery().Must()

	if queries.QueryField != "" && queries.FieldValue != "" {
		query.Must(elastic.NewMatchQuery(fmt.Sprint(queries.QueryField), queries.FieldValue))
	}

	from, _ := strconv.Atoi(queries.OffSet)
	size, _ := strconv.Atoi(queries.Limit)

	res, err := repo.elastic.
		Search("company").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllCompanies", "Cant get all companies").Error("Error - " + err.Error())
		return []models.Company{}, err
	}

	var ttyp models.Company
	companies := make([]models.Company, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		company := item.(models.Company)
		companies = append(companies, company)
	}

	return companies, nil
}

func (repo *companyRepository) SearchCompany(searchText, offSet, limit string) ([]models.Company, error) {
	query := elastic.NewBoolQuery().Must()

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"legal_name",
				"address",
				"city",
				"country",
				"phone_number",
				"fax_number",
				"billing_address",
				"billing_method",
				"billing_type",
				"billing_email",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("company").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchCompany", "Cant get all company by search").Error("Error - " + err.Error())
		return []models.Company{}, err
	}

	var ttyp models.Company
	companies := make([]models.Company, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		company := item.(models.Company)
		companies = append(companies, company)
	}

	return companies, nil
}

func (repo *companyRepository) GetAllCompaniesForReport() ([]models.Company, error) {

	res, err := repo.elastic.
		Search("company").
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllCompaniesForReport", "Cant get all companies for report").Error("Error - " + err.Error())
		return []models.Company{}, err
	}

	var ttyp models.Company
	companies := make([]models.Company, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		company := item.(models.Company)
		companies = append(companies, company)
	}

	return companies, nil
}
