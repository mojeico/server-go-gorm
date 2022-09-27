package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nats-io/nats.go"

	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/queue"
	"github.com/trucktrace/pkg/reports"

	"github.com/gin-gonic/gin"
)

type SettlementController interface {
	CreateSettlement(ctx *gin.Context)
	DeleteSettlement(ctx *gin.Context)
	UpdateSettlement(ctx *gin.Context)
	GetAllSettlement(ctx *gin.Context)
	GetSettlementById(ctx *gin.Context)
	SearchSettlements(ctx *gin.Context)
	CreateAllSettlementsReport(ctx *gin.Context)
	GetSettlementReportForDriverPDF(ctx *gin.Context)
}

type settlementController struct {
	charge          service.ChargesService
	order           service.OrderService
	service         service.SettlementService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewSettlementController(setService service.SettlementService, redisConnection *redis.Client, order service.OrderService, charge service.ChargesService, nats *nats.Conn) SettlementController {
	return &settlementController{
		service:         setService,
		redisConnection: redisConnection,
		order:           order,
		charge:          charge,
		nats:            nats,
	}
}

func (controller *settlementController) CreateSettlement(ctx *gin.Context) {

	var queries = map[string]string{
		"id": ctx.Param("orderId"),
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_settlement", Body: nil, Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateSettlement", "Can't create Settlement").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't create Settlement", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, "OK")

}

func (controller *settlementController) DeleteSettlement(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{

		"id":         ctx.Param("id"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	message := queue.Message{Method: "delete_settlement", Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteSettlement", "Can't delete Settlement").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Settlement", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Settlement was deleted").Info(fmt.Sprintf("Settlement was deleted with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *settlementController) UpdateSettlement(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var settlement models.SettlementUpdateInput
	if err := ctx.BindJSON(&settlement); err != nil {
		logger.ErrorLogger("UpdateSettlement", "Can't bind SettlementUpdateInput json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding SettlementUpdateInput json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{

		"id":         ctx.Param("id"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	message := queue.Message{Method: "update_settlement", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateSettlement", "Can't update Settlement").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating Settlement", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Settlement was updated").Info(fmt.Sprintf("Settlement was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *settlementController) GetAllSettlement(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	settlements, err := controller.service.GetAllSettlements(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get all settlements", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all settlements", settlements)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllSettlement").Info("Return all settlements")

}

func (controller *settlementController) GetSettlementById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	settlement, err := controller.service.GetSettlementById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get settlement by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return one settlement", settlement)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetSettlementById").Info("Return one settlement")

}

func (controller *settlementController) SearchSettlements(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	settlements, err := controller.service.SearchSettlements(searchText, offSet, limit)

	if err != nil {
		res := response.BuildErrorResponse("can't get all settlements by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all settlements  by search", settlements)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchSettlements").Info("Return all settlements  by search")

}

func (controller *settlementController) CreateAllSettlementsReport(ctx *gin.Context) {

	settlements, err := controller.service.GetAllSettlementsForReport()

	if err != nil {
		res := response.BuildErrorResponse("can't get all settlements for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	settlementModel := models.Settlement{}
	header := reports.GetHeaderAsArray(settlementModel)

	var data [][]string
	for _, v := range settlements {
		settlement := reports.GetValuesAsArray(v)
		data = append(data, settlement)
	}

	pdf := reports.NewReport(870, "Settlement")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllSettlementsReport", "can't create get all settlements report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all settlements report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	if err = pdf.Output(ctx.Writer); err != nil {
		logger.ErrorLogger("CreateAllSettlementsReport", "can't save get all settlements report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all settlements report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()
}

func (controller *settlementController) GetSettlementReportForDriverPDF(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	settlement, err := controller.service.GetSettlementById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get settlement by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	queryParams = helper.NewOneQueryParams(
		fmt.Sprint(settlement.OrderId), "Completed", "false", "false",
	)
	order, err := controller.order.GetOrderById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get order for creating settlement by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var chargesList []models.Charges
	for _, chargeId := range order.ChargesList {
		queryParams = helper.NewOneQueryParams(
			string(chargeId), "Completed", "false", "false",
		)
		charge, err := controller.charge.GetChargeById(queryParams)
		if err != nil {
			res := response.BuildErrorResponse("problems with finding charges", err.Error(), response.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		chargesList = append(chargesList, charge)

	}

	pdf := reports.CreateSettlementReport(settlement, order, chargesList)

	settlementBytes, err := pdf.Output()

	if err != nil {
		logger.ErrorLogger("GetSettlementReportForDriverPDF", "can't generate settlement pdf file").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't generate settlement pdf file", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := response.BuildResponse("Document was created", settlementBytes)
	ctx.JSON(http.StatusOK, res)

}
