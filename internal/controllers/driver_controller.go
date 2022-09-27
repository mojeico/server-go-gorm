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

type DriverController interface {
	CreateDriver(ctx *gin.Context)
	GetDriverById(ctx *gin.Context)
	GetAllDrivers(ctx *gin.Context)
	GetAllDriversByCompanyId(ctx *gin.Context)
	DeleteDriver(ctx *gin.Context)
	UpdateDriver(ctx *gin.Context)
	SearchDrivers(ctx *gin.Context)
	CreateAllDriversReport(ctx *gin.Context)
}

type driverController struct {
	service         service.DriverService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewDriverController(service service.DriverService, redisConnection *redis.Client, nats *nats.Conn) DriverController {
	return &driverController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller *driverController) CreateDriver(ctx *gin.Context) {

	var driver models.Driver
	if err := ctx.BindJSON(&driver); err != nil {
		logger.ErrorLogger("CreateDriver", "Can't bind Driver json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Driver json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID
	driver.CompanyId = companyId

	jsonData, _ := json.Marshal(driver)
	err := validation.ValidateJSON(jsonData, driver)
	if err != nil {
		logger.ErrorLogger("CreateDriver", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_driver", Body: jsonData, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateDriver", "Problem with creating Driver").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating Driver", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Driver was created").Info(fmt.Sprintf("Driver was created with first name %s.", driver.FirstName))
	ctx.JSON(http.StatusOK, driver)
}

func (controller *driverController) GetDriverById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	driver, err := controller.service.GetDriverById(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get a driver by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("gotten driver successfully", driver)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetDriverById").Info(fmt.Sprintf("Driver was returned with id - %T", driver.ID))

}

func (controller *driverController) GetAllDrivers(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	drivers, err := controller.service.GetAllDrivers(queryParams, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all drivers", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("gotten all drivers successfully", drivers)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllDrivers").Info("Return all drivers")

}

func (controller *driverController) GetAllDriversByCompanyId(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	companyId := ctx.Param("companyId")
	orderBy := ctx.Query("orderBy")
	orderDir, _ := strconv.Atoi(ctx.Query("orderDir"))

	drivers, err := controller.service.GetAllDriversByCompanyId(companyId, queryParams, orderBy, byte(orderDir))
	if err != nil {
		res := response.BuildErrorResponse("can't get all drivers", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("gotten all drivers successfully", drivers)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllDriversByCompanyId").Info("Return driver list company id")

}

func (controller *driverController) DeleteDriver(ctx *gin.Context) {

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

	message := queue.Message{Method: "delete_driver", Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteDriver", "Can't delete Driver").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Driver", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	logger.InfoLogger("Customer was deleted").Info(fmt.Sprintf("Customer was deleted with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *driverController) UpdateDriver(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
	var driver models.DriverUpdateInput

	if err := ctx.BindJSON(&driver); err != nil {
		logger.ErrorLogger("UpdateDriver", "Can't bind DriverUpdateInput json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding DriverUpdateInput json", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_driver", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateDriver", "Can't update Driver").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't update Driver", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	logger.InfoLogger("Customer was updated").Info(fmt.Sprintf("Customer was updated with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")

}

func (controller *driverController) SearchDrivers(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	drivers, err := controller.service.SearchDrivers(searchText, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get drivers by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("gotten drivers by search successfully", drivers)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchDrivers").Info(fmt.Sprintf("Return all drivers by search. searchText - %s", searchText))

}

func (controller *driverController) CreateAllDriversReport(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	drivers, err := controller.service.GetAllDriversForReport(companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all drivers for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	driverModel := models.Driver{}
	header := reports.GetHeaderAsArray(driverModel)

	var data [][]string
	for _, v := range drivers {
		driver := reports.GetValuesAsArray(v)
		data = append(data, driver)
	}

	pdf := reports.NewReport(2670, "Driver")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllDriversReport", "can't create get all drivers report").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't create get all drivers report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	err = pdf.Output(ctx.Writer)

	if err != nil {
		logger.ErrorLogger("CreateAllDriversReport", "can't save get all drivers report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all drivers report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()

	logger.ErrorLogger("CreateAllDriversReport", "Report was created").Error("Error - " + err.Error())
}
