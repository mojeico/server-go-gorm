package repository

import (
	"github.com/olivere/elastic/v7"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/pkg/logger"
	"reflect"
)

type TrailerCommentRepository interface {
	GetTrailerCommentsByTrailerID(trailerId string) ([]models.TrailerComment, error)
}

type trailerCommentRepository struct {
	elastic *elastic.Client
}

func NewTrailerCommentRepository(elastic *elastic.Client) TrailerCommentRepository {
	return &trailerCommentRepository{
		elastic: elastic,
	}
}

func (repo trailerCommentRepository) GetTrailerCommentsByTrailerID(trailerId string) ([]models.TrailerComment, error) {

	query := elastic.NewBoolQuery().Must()
	query.Must(elastic.NewMatchQuery("trailer_id", trailerId))

	res, err := repo.elastic.
		Search("trailer_comments").
		Query(query).
		Do(ctx)

	if err != nil {
		logger.ErrorLogger("GetTrailerCommentsByTrailerID", "Cant get all comments by trailer id").Error("Error - " + err.Error())
		return nil, err
	}

	var ttyp models.TrailerComment
	comments := make([]models.TrailerComment, 0)

	for _, item := range res.Each(reflect.TypeOf(ttyp)) {
		comment := item.(models.TrailerComment)
		comments = append(comments, comment)

	}

	return comments, nil
}
