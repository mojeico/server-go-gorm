package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/queue"
	"github.com/trucktrace/pkg/validation"
)

type TrailerCommentController interface {
	GetTrailerCommentsByTrailerID(ctx *gin.Context)
	CreateTrailerComment(ctx *gin.Context)
}

type trailerCommentController struct {
	service         service.TrailerCommentService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewTrailerCommentController(service service.TrailerCommentService, redisConnection *redis.Client, nats *nats.Conn) TrailerCommentController {
	return &trailerCommentController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller trailerCommentController) GetTrailerCommentsByTrailerID(ctx *gin.Context) {

	trailerId := ctx.Param("trailerId")

	comments, err := controller.service.GetTrailerCommentsByTrailerID(trailerId)

	if err != nil {
		res := response.BuildErrorResponse("can't get trailers comments", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Trailer comments  was gotten", comments)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetTrailerCommentsByTrailerID").Info("can't get trailers comments")

}

func (controller trailerCommentController) CreateTrailerComment(ctx *gin.Context) {
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	err := validation.ValidateJSON(jsonData, models.TrailerComment{})
	if err != nil {
		logger.ErrorLogger("CreateTrailerComment", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_trailer_comments", Body: jsonData, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateTrailerComment", "Can't create trailer comment").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating trailer comment", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Trailer comments was created").Info("Trailer comment was created ")
	ctx.JSON(http.StatusOK, "OK")
}
