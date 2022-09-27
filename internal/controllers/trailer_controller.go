package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/queue"
	"github.com/trucktrace/pkg/reports"
)

type TrailerController interface {
	GetTrailerByID(ctx *gin.Context)
	GetAllTrailersByCompanyId(ctx *gin.Context)
	GetAllTrailers(ctx *gin.Context)
	CreateTrailer(ctx *gin.Context)
	DeleteTrailer(ctx *gin.Context)
	UpdateTrailer(ctx *gin.Context)
	SearchTrailers(ctx *gin.Context)
	CreateAllTrailersReport(ctx *gin.Context)
}

type trailerController struct {
	service         service.TrailerService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewTrailerController(service service.TrailerService, redisConnection *redis.Client, nats *nats.Conn) TrailerController {
	return &trailerController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller *trailerController) GetTrailerByID(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	trailer, err := controller.service.GetTrailerById(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get trailer by id", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("trailer by id was gotten", trailer)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetTrailerByID").Info("trailer by id was gotten")

}

func (controller *trailerController) GetAllTrailers(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	trailers, err := controller.service.GetAllTrailers(queryParams, companyId)
	if err != nil {
		res := response.BuildErrorResponse("can't get all trailers", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("trailers was gotten", trailers)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllTrailers").Info("trailers was gotten")
}

func (controller *trailerController) GetAllTrailersByCompanyId(ctx *gin.Context) {

	companyId := ctx.Param("companyId")

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	orderBy := ctx.Query("orderBy")
	orderDir, _ := strconv.Atoi(ctx.Query("orderDir"))

	trailer, err := controller.service.GetAllTrailerByCompanyId(companyId, queryParams, orderBy, byte(orderDir))

	if err != nil {
		res := response.BuildErrorResponse("can't get all trailer", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("trailer was gotten", trailer)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllTrailersByCompanyId").Info("trailer was gotten")

}

func (controller *trailerController) CreateTrailer(ctx *gin.Context) {

	var trailer models.Trailer

	if err := ctx.ShouldBindJSON(&trailer); err != nil {
		logger.ErrorLogger("CreateTrailer", "Can't bind Trailer json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Trailer json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	companyId := userContext.(models.User).CompanyID
	trailer.CompanyId = companyId
	jsonData, _ := json.Marshal(trailer)

	err := validation.ValidateJSON(jsonData, trailer)
	if err != nil {
		logger.ErrorLogger("CreateTrailer", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	message := queue.Message{Method: "create_trailer", Body: jsonData, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateTrailer", "Problem with creating trailer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating trailer", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Trailer was created").Info(fmt.Sprintf("Trailer was created with name %s.", trailer.Name))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *trailerController) DeleteTrailer(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email
	companyId := userContext.(models.User).CompanyID

	var queries = map[string]string{

		"id":         ctx.Param("id"),
		"company_id": strconv.Itoa(companyId),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	message := queue.Message{Method: "delete_trailer", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteTrailer", "Can't delete Trailer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Trailer", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Trailer was deleted").Info(fmt.Sprintf("Trailer was deleted with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *trailerController) UpdateTrailer(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var trailer models.TrailerInputUpdate
	if err := ctx.BindJSON(&trailer); err != nil {
		logger.ErrorLogger("UpdateTrailer", "Can't bind TrailerInputUpdate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding TrailerInputUpdate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	companyId := userContext.(models.User).CompanyID

	var queries = map[string]string{

		"id":         ctx.Param("id"),
		"company_id": strconv.Itoa(companyId),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	message := queue.Message{Method: "update_trailer", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateTrailer", "Can't update Trailer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating Trailer", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Trailer was updated").Info(fmt.Sprintf("Trailer was updated with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")

}

func (controller *trailerController) SearchTrailers(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	trailers, err := controller.service.SearchTrailers(searchText, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all trailers by search", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("trailers was gotten by search", trailers)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchTrailers").Info("trailers was gotten by searchn")

}

func (controller *trailerController) CreateAllTrailersReport(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	trailers, err := controller.service.GetAllTrailersForReport(companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all trailer for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	trailerModel := models.Trailer{}
	header := reports.GetHeaderAsArray(trailerModel)

	var data [][]string
	for _, v := range trailers {
		trailer := reports.GetValuesAsArray(v)
		data = append(data, trailer)
	}

	pdf := reports.NewReport(2070, "Trailer")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllTrailersReport", "can't create get all trailer report").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't create get all trailer report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	if err = pdf.Output(ctx.Writer); err != nil {
		logger.ErrorLogger("CreateAllTrailersReport", "can't save get all trailer report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all trailer report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()

	logger.ErrorLogger("CreateAllTrailersReport", "Report was created").Error("Error - " + err.Error())
}
