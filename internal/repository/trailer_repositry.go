package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/trucktrace/pkg/logger"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/helper"
)

type TrailerRepository interface {
	GetAllTrailerByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]map[string]interface{}, error)
	GetTrailerById(queries *helper.OneQueryParams) (map[string]interface{}, error)
	GetAllTrailers(queries *helper.GetAllQueryParams, companyId int) ([]map[string]interface{}, error)
	SearchTrailers(searchText, offSet, limit string, companyId int) ([]map[string]interface{}, error)
	GetAllTrailersForReport(companyId int) ([]models.Trailer, error)
}

type trailerRepository struct {
	pg       *gorm.DB
	elastic  *elastic.Client
	comments TrailerCommentRepository
}

func NewTrailerRepository(pg *gorm.DB, elastic *elastic.Client, comments TrailerCommentRepository) TrailerRepository {
	return &trailerRepository{
		pg:       pg,
		elastic:  elastic,
		comments: comments,
	}
}

func (repo *trailerRepository) GetAllTrailers(queries *helper.GetAllQueryParams, companyId int) ([]map[string]interface{}, error) {

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
		Search("trailer").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllTrailers", "Cant get all trailers").Error("Error - " + err.Error())
		return nil, err
	}

	var trailers []map[string]interface{}

	for i := 0; i < int(res.TotalHits()); i++ {
		var trailer map[string]interface{}
		_ = json.Unmarshal(res.Hits.Hits[i].Source, &trailer)
		id := fmt.Sprint(trailer["ID"])
		comment, _ := repo.comments.GetTrailerCommentsByTrailerID(id)
		trailer["comments"] = comment
		trailers = append(trailers, trailer)
	}

	return trailers, nil
}

func (repo *trailerRepository) GetAllTrailerByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]map[string]interface{}, error) {

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
		Search("trailer").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllTrailerByCompanyId", "Cant get all trailers by company id").Error("Error - " + err.Error())
		return nil, err
	}

	var trailers []map[string]interface{}

	for i := 0; i < int(res.TotalHits()); i++ {
		var trailer map[string]interface{}
		_ = json.Unmarshal(res.Hits.Hits[i].Source, &trailer)
		id := fmt.Sprint(trailer["ID"])
		comment, _ := repo.comments.GetTrailerCommentsByTrailerID(id)
		trailer["comments"] = comment
		trailers = append(trailers, trailer)
	}

	return trailers, nil
}

func (repo *trailerRepository) GetTrailerById(queries *helper.OneQueryParams) (map[string]interface{}, error) {

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
		Search("trailer").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetTrailerById", "Cant get trailer by id").Error("Error - " + err.Error())
		return nil, err
	}

	var trailer map[string]interface{}

	if res.TotalHits() > 0 {
		_ = json.Unmarshal(res.Hits.Hits[0].Source, &trailer)
		id := fmt.Sprint(trailer["ID"])
		comment, _ := repo.comments.GetTrailerCommentsByTrailerID(id)
		trailer["comments"] = comment
		return trailer, nil
	}

	return nil, errors.New("trailer not found")
}

func (repo *trailerRepository) SearchTrailers(searchText, offSet, limit string, companyId int) ([]map[string]interface{}, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"name",
				"unit_number",
				"make",
				"plate",
				"state",
				"vin_number",
				"owner_name",
				"location",
				"status",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("trailer").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchTrailers", "Cant get all trailers by search").Error("Error - " + err.Error())
		return nil, err
	}

	var trailers []map[string]interface{}

	for i := 0; i < int(res.TotalHits()); i++ {
		var trailer map[string]interface{}
		_ = json.Unmarshal(res.Hits.Hits[i].Source, &trailer)
		id := fmt.Sprint(trailer["ID"])
		comment, _ := repo.comments.GetTrailerCommentsByTrailerID(id)
		trailer["comments"] = comment
		trailers = append(trailers, trailer)
	}
	return trailers, nil
}

func (repo *trailerRepository) GetAllTrailersForReport(companyId int) ([]models.Trailer, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	res, err := repo.elastic.
		Search("trailer").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllTrailersForReport", "Cant get all trailers for report").Error("Error - " + err.Error())
		return []models.Trailer{}, err
	}

	var ttyp models.Trailer
	trailers := make([]models.Trailer, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		trailer := item.(models.Trailer)
		trailers = append(trailers, trailer)
	}

	return trailers, nil
}
