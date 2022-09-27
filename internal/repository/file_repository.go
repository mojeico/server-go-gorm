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

type FileRepository interface {
	GetFileById(queries *helper.OneQueryParams) (models.File, error)
	GetLastFileByOwnerId(ownerType string, queries *helper.OneQueryParams) (models.File, error)
	GetAllFilesByOwnerId(ownerType string, queries *helper.GetAllQueryParams) ([]models.File, error)
	GetAllFiles(queries *helper.GetAllQueryParams) ([]models.File, error)
	SearchFiles(searchText, offSet, limit string) ([]models.File, error)
}

type fileRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewFileRepository(pg *gorm.DB, elastic *elastic.Client) FileRepository {
	return &fileRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *fileRepository) GetFileById(queries *helper.OneQueryParams) (models.File, error) {

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
		Search("files").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetFileById", "Cant get file by id ").Error("Error - " + err.Error())
		return models.File{}, err
	}

	var ttyp models.File

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		file := item.(models.File)
		return file, nil
	}

	return models.File{}, errors.New("file not found")
}

func (repo *fileRepository) GetLastFileByOwnerId(ownerType string, queries *helper.OneQueryParams) (models.File, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("_id", queries.Id))
	query.Must(elastic.NewMatchQuery("owner_type", ownerType))

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
		Search("files").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetLastFileByOwnerId", "Cant get last file by owner id").Error("Error - " + err.Error())
		return models.File{}, err
	}

	var ttyp models.File

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		file := item.(models.File)
		return file, nil
	}

	return models.File{}, errors.New("file not found")
}

func (repo *fileRepository) GetAllFilesByOwnerId(ownerType string, queries *helper.GetAllQueryParams) ([]models.File, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("owner_type", ownerType))

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
		Search("files").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllFilesByOwnerId", "Cant get all files by owner id").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.File
	files := make([]models.File, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		file := item.(models.File)
		files = append(files, file)
	}

	return files, nil
}

func (repo *fileRepository) GetAllFiles(queries *helper.GetAllQueryParams) ([]models.File, error) {

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
		Search("files").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllFiles", "Cant get all files").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.File
	files := make([]models.File, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		file := item.(models.File)
		files = append(files, file)
	}

	return files, nil
}

func (repo *fileRepository) SearchFiles(searchText, offSet, limit string) ([]models.File, error) {

	query := elastic.NewBoolQuery().Must()

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"name",
				"extension",
				"expiration_status",
				"comment",
				"owner_type",
				"status",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("files").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchFiles", "Cant get all files by search").Error("Error - " + err.Error())
		return []models.File{}, err
	}

	var ttyp models.File
	files := make([]models.File, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		file := item.(models.File)
		files = append(files, file)
	}

	return files, nil
}
