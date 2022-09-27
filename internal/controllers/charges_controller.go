package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

type ChargesController interface {
	CreateCharges(ctx *gin.Context)
	UpdateCharges(ctx *gin.Context)
	DeleteCharges(ctx *gin.Context)

	GetAllCharges(ctx *gin.Context)
	GetAllChargesBySettlementId(ctx *gin.Context)
	GetAllChargesByOrderId(ctx *gin.Context)
	GetChargeById(ctx *gin.Context)

	SearchCharge(ctx *gin.Context)

	CreateAllChargersReport(ctx *gin.Context)
}

type chargerController struct {
	service         service.ChargesService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewChargesController(service service.ChargesService, redisConnection *redis.Client, nats *nats.Conn) ChargesController {
	return &chargerController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller *chargerController) CreateCharges(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)

	var charges models.Charges
	err := validation.ValidateJSON(jsonData, charges)
	if err != nil {
		logger.ErrorLogger("CreateCharges", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := json.Unmarshal(jsonData, &charges); err != nil {
		logger.ErrorLogger("CreateCharges", "Can't Unmarshal Charges json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with Unmarshal Charges json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := helper.CheckChargesRate(charges); err != nil {
		errorMessage := fmt.Sprintf("Inserting wrong Rate - %v and Type - %v of Charges ", charges.Rate, charges.TypeDeductions)
		logger.ErrorLogger("CreateCharges", errorMessage).Error("Error - " + err.Error())
		res := response.BuildErrorResponse(errorMessage, err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_charges", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateCharges", "Can't create charges").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating charges", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Charges was created").Info(fmt.Sprintf("Charges was created for dririver %s and company %s.", charges.DriverName, charges.CompanyName))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *chargerController) UpdateCharges(ctx *gin.Context) {

	var charges models.ChargesUpdateInput
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	if err := ctx.BindJSON(&charges); err != nil {
		logger.ErrorLogger("UpdateCharges", "Can't bind ChargesUpdateInput json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating Charges", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_charges", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateCharges", "Can't update charges").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't update charges", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Charges was updated").Info(fmt.Sprintf("Charges was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *chargerController) DeleteCharges(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{
		"id":         ctx.Param("id"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	message := queue.Message{Method: "delete_charges", Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteCharges", "Can't delete charges").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't delete charges", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Charges was deleted").Info(fmt.Sprintf("Charges was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *chargerController) GetAllCharges(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	charges, err := controller.service.GetAllCharges(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get all Charges", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all charges", charges)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllCharges").Info("Charges list was returned")
}

func (controller *chargerController) GetAllChargesBySettlementId(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParamsWithId(
		ctx.Param("settlementId"),
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	customer, err := controller.service.GetAllChargesBySettlementId(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get changes by settlement id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return all changes by settlement id", customer)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllChargesBySettlementId").Info("Charges list by settlement id  was returned")
}

func (controller *chargerController) GetAllChargesByOrderId(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParamsWithId(
		ctx.Param("orderId"),
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	customer, err := controller.service.GetAllChargesByOrderId(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get changes by order id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return all changes by order id", customer)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllChargesByOrderId").Info("Charges list by order id was returned")

}

func (controller *chargerController) GetChargeById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	charge, err := controller.service.GetChargeById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get charges by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return one charges", charge)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetChargeById").Info(fmt.Sprintf("Charge with id - %T was returned", charge.ID))

}

func (controller chargerController) SearchCharge(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	charges, err := controller.service.SearchCharge(searchText, offSet, limit)

	if err != nil {
		res := response.BuildErrorResponse("can't get all charges by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all charges by search", charges)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetChargeById").Info(fmt.Sprintf("Charges list by search  was returned - searchText = %s ", searchText))

}

func (controller *chargerController) CreateAllChargersReport(ctx *gin.Context) {

	charges, err := controller.service.GetAllChargesForReport()

	if err != nil {
		res := response.BuildErrorResponse("can't get all charges for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	chargeModel := models.Charges{}
	header := reports.GetHeaderAsArray(chargeModel)

	var data [][]string
	for _, v := range charges {
		charge := reports.GetValuesAsArray(v)
		data = append(data, charge)
	}

	pdf := reports.NewReport(870, "Charges")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllChargersReport", "Can't create get all charges report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all charges report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	//err = reports.SavePDF(pdf)
	err = pdf.Output(ctx.Writer)

	if err != nil {
		logger.ErrorLogger("CreateAllChargersReport", "can't save get all charges report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all charges report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()

	logger.InfoLogger("GetChargeById").Info("All chargers report was created")

}
