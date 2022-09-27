package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nats-io/nats.go"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/queue"
	"github.com/trucktrace/pkg/reports"
	"github.com/trucktrace/pkg/validation"
)

type SafetyController interface {
	CreateSafety(ctx *gin.Context)
	GetSafetyById(ctx *gin.Context)
	GetAllSafeties(ctx *gin.Context)
	GetAllSafetiesByCompanyId(ctx *gin.Context)
	DeleteSafety(ctx *gin.Context)
	UpdateSafety(ctx *gin.Context)
	SearchSafety(ctx *gin.Context)
	GetAllSafetiesByType(ctx *gin.Context)
	CreateAllSafetyReport(ctx *gin.Context)
}

type safetyController struct {
	service         service.SafetyService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewSafetyController(service service.SafetyService, redisConnection *redis.Client, nats *nats.Conn) SafetyController {
	return &safetyController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller *safetyController) CreateSafety(ctx *gin.Context) {

	var safety models.Safety
	if err := ctx.BindJSON(&safety); err != nil {
		logger.ErrorLogger("CreateSafety", "Can't bind Safety json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Safety json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID
	safety.CompanyID = companyId

	jsonData, _ := json.Marshal(safety)

	err := validation.ValidateJSON(jsonData, safety)
	if err != nil {
		logger.ErrorLogger("CreateSafety", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_safety", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateSafety", "Can't create Safety").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating Safety", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Safety was created").Info(fmt.Sprintf("Safety was created with type %s.", safety.SafetyType))

	ctx.JSON(http.StatusOK, "OK")

}

func (controller *safetyController) GetSafetyById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)
	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	safety, err := controller.service.GetSafetyById(strconv.Itoa(companyId), queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get safety by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return one safety", safety)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetSafetyById").Info("Return one safety by id")

}

func (controller *safetyController) GetAllSafeties(ctx *gin.Context) {

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

	safeties, err := controller.service.GetAllSafeties(queryParams, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all safeties", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all safeties", safeties)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllSafeties").Info("Return all safeties")

}

func (controller *safetyController) GetAllSafetiesByCompanyId(ctx *gin.Context) {

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

	safeties, err := controller.service.GetAllSafetiesByCompanyId(companyId, queryParams, orderBy, byte(orderDir))

	if err != nil {
		res := response.BuildErrorResponse("Can't get all safeties by company id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Get all safeties by company id successfully", safeties)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllSafetiesByCompanyId").Info("Get all safeties by company id successfully")

}

func (controller *safetyController) UpdateSafety(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var safety models.SafetyInputUpdate
	if err := ctx.BindJSON(&safety); err != nil {
		logger.ErrorLogger("UpdateSafety", "Can't bind SafetyInputUpdate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding SafetyInputUpdate json", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_safety", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateSafety", "Can't update  Safety").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating Safety", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Safety was updated").Info(fmt.Sprintf("Safety was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *safetyController) DeleteSafety(ctx *gin.Context) {
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

	message := queue.Message{Method: "delete_safety", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteSafety", "Can't delete  Safety").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Safety", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Safety was deleted").Info(fmt.Sprintf("Safety was deleted with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *safetyController) SearchSafety(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	safeties, err := controller.service.SearchSafety(searchText, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all safeties by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all safeties by search", safeties)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchSafety").Info("Return all safeties by search")

}

func (controller *safetyController) GetAllSafetiesByType(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")

	safetyType := ctx.Param("type")

	if safetyType != "reports" && safetyType != "permits" && safetyType != "registration" {
		logger.ErrorLogger("UpdateSafety", "Wrong safety type").Error(fmt.Sprintf("Wrong safety type - %s ", safetyType))
		res := response.BuildErrorResponse(fmt.Sprintf("Wrong safety type - %s ", safetyType), "Wrong safety type", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	safetiesByType, err := controller.service.GetAllSafetyByType(safetyType, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all safeties by type", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all safeties by type", safetiesByType)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllSafetiesByType").Info("Return all safeties by type")

}

func (controller *safetyController) CreateAllSafetyReport(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	safeties, err := controller.service.GetAllSafetiesForReport(companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all safeties for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	safetyModel := models.Safety{}
	header := reports.GetHeaderAsArray(safetyModel)

	var data [][]string
	for _, v := range safeties {
		charge := reports.GetValuesAsArray(v)
		data = append(data, charge)
	}

	pdf := reports.NewReport(710, "Safety")
	pdf = reports.Header(pdf, header, 40)
	pdf = reports.Table(pdf, data, 40)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllSafetyReport", "Can't create get all safeties report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all safeties report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	//err = reports.SavePDF(pdf)

	if err = pdf.Output(ctx.Writer); err != nil {
		logger.ErrorLogger("CreateAllSafetyReport", "can't save get all safeties report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all safeties report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	pdf.Close()
}
