package repository

import (
	"encoding/json"
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

type UserRepository interface {
	GetAllUsers(queries *helper.GetAllQueryParams, companyId int) ([]models.User, error)
	GetAllUsersByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.User, error)
	GetUserById(queries *helper.OneQueryParams) (models.User, error)
	GetUserByIdForMiddleware(id string) (models.User, error)
	SearchUsers(searchText, offSet, limit string, companyId int) ([]models.User, error)
}

type userRepository struct {
	pg      *gorm.DB
	elastic *elastic.Client
}

func NewUserRepository(pg *gorm.DB, elastic *elastic.Client) UserRepository {
	return &userRepository{
		pg:      pg,
		elastic: elastic,
	}
}

func (repo *userRepository) GetAllUsers(queries *helper.GetAllQueryParams, companyId int) ([]models.User, error) {

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
		Search("user").
		Query(query).
		//Sort("email.keyword", false).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllUsers", "Cant get all users").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.User
	users := make([]models.User, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		user := item.(models.User)
		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) GetAllUsersByCompanyId(companyId string, queries *helper.GetAllQueryParams) ([]models.User, error) {

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
		Search("user").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetAllUsersByCompanyId", "Cant get all users by company id").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.User
	users := make([]models.User, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		user := item.(models.User)
		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) GetUserById(queries *helper.OneQueryParams) (models.User, error) {

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
		Search("user").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetUserById", "Cant get user by id").Error("Error - " + err.Error())
		return models.User{}, err
	}

	var ttyp models.User

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		user := item.(models.User)
		return user, nil
	}

	return models.User{}, errors.New("user not found")
}

func (repo *userRepository) GetUserByIdForMiddleware(userId string) (models.User, error) {

	userResponse, err := repo.elastic.
		Get().
		Index("user").
		Id(userId).
		Do(ctx)

	if err != nil {
		logger.SystemLoggerError("GetUserByIdForMiddleware", "Cant get user by id for middleware").Error("Error - " + err.Error())
		return models.User{}, err
	}
	if userResponse.Found {
		var user models.User
		err := json.Unmarshal(userResponse.Source, &user)
		if err != nil {
			return models.User{}, err

		}
		return user, nil
	}
	return models.User{}, errors.New("user not found")

}

func (repo *userRepository) SearchUsers(searchText, offSet, limit string, companyId int) ([]models.User, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("company_id", companyId))

	if searchText != "" {
		query.Must(
			elastic.NewMultiMatchQuery(searchText,
				"name",
				"email",
				"phone",
				"username",
				"status",
				"role",
			).Type("phrase_prefix"))
	}

	from, _ := strconv.Atoi(offSet)
	size, _ := strconv.Atoi(limit)

	res, err := repo.elastic.
		Search("user").
		Query(query).
		From(from).
		Size(size).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("SearchUsers", "Cant get all users by search").Error("Error - " + err.Error())
		return []models.User{}, err
	}

	var ttyp models.User
	users := make([]models.User, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		user := item.(models.User)
		users = append(users, user)
	}

	return users, nil
}
