package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nats-io/nats.go"

	"github.com/trucktrace/pkg/reports"
	"github.com/trucktrace/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/queue"
)

var emptyObj = response.EmptyObj{}

type TruckController interface {
	CreateTruck(ctx *gin.Context)
	GetAllTrucks(ctx *gin.Context)
	GetAllTrucksByCompanyId(ctx *gin.Context)
	GetTruckById(ctx *gin.Context)
	UpdateTruck(ctx *gin.Context)
	DeleteTruck(ctx *gin.Context)
	SearchTrucks(ctx *gin.Context)
	CreateAllTrucksReport(ctx *gin.Context)
}

type truckController struct {
	service         service.TruckService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewTruckController(service service.TruckService, redisConnection *redis.Client, nats *nats.Conn) TruckController {
	return &truckController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

//-----------------------+CREATE+-----------------------//
func (controller *truckController) CreateTruck(ctx *gin.Context) {

	var truck models.Truck
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	if err := json.Unmarshal(jsonData, &truck); err != nil {
		logger.ErrorLogger("CreateTruck", "Can't bind Truck json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Truck json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	companyId := userContext.(models.User).CompanyID
	truck.CompanyId = companyId
	jsonData, _ = json.Marshal(truck)
	err := validation.ValidateJSON(jsonData, truck)
	if err != nil {
		logger.ErrorLogger("CreateTruck", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	message := queue.Message{Method: "create_truck", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateTruck", "Can't create Truck").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't create Truck", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Truck was created").Info(fmt.Sprintf("Truck was created with plate %T.", truck.Plate))
	ctx.JSON(http.StatusOK, "OK")
}

//-----------------------+GET ALL TRUCKS+-----------------------//
func (controller *truckController) GetAllTrucks(ctx *gin.Context) {

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

	trucks, err := controller.service.GetAllTrucks(queryParams, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all trucks", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("trucks was gotten", trucks)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllTrucks").Info("trucks was gotten")
}

//-----------------------+GET TRUCK BY COMPANY ID+-----------------------//
func (controller *truckController) GetAllTrucksByCompanyId(ctx *gin.Context) {

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

	trucks, err := controller.service.GetAllTrucksByCompanyId(companyId, queryParams, orderBy, byte(orderDir))

	if err != nil {
		res := response.BuildErrorResponse("can't get all trucks", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("trucks was gotten", trucks)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllTrucksByCompanyId").Info("trucks was gotten")
}

//-----------------------+GET TRUCK BY ID+-----------------------//
func (controller *truckController) GetTruckById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	truck, err := controller.service.GetTruckById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get truck by id", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("truck by id was gotten", truck)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetTruckById").Info("truck by id was gotten")

}

// -----------------------+UPDATE+-----------------------//
func (controller *truckController) UpdateTruck(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var truck models.TruckUpdateInput

	if err := ctx.BindJSON(&truck); err != nil {
		logger.ErrorLogger("UpdateTruck", "Can't bind TruckUpdateInput json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding TruckUpdateInput json", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_truck", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateTruck", "Can't update Truck").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating Truck", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Truck was updated").Info(fmt.Sprintf("Truck was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

//-----------------------+DELETE+-----------------------//
func (controller *truckController) DeleteTruck(ctx *gin.Context) {

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

	message := queue.Message{Method: "delete_truck", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteTruck", "Can't delete Truck").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Truck", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Truck was deleted").Info(fmt.Sprintf("Truck was deleted with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")

}

func (controller *truckController) SearchTrucks(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	trucks, err := controller.service.SearchTrucks(searchText, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all trucks by search", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("trucks was gotten by search", trucks)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchTrucks").Info("can't get all trucks by search")

}

func (controller *truckController) CreateAllTrucksReport(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	trucks, err := controller.service.GetAllTrucksForReport(companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all trucks for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	truckModel := models.Truck{}
	header := reports.GetHeaderAsArray(truckModel)

	var data [][]string
	for _, v := range trucks {
		truck := reports.GetValuesAsArray(v)
		data = append(data, truck)
	}

	pdf := reports.NewReport(1850, "Truck")
	pdf = reports.Header(pdf, header, 70)
	pdf = reports.Table(pdf, data, 70)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllTrucksReport", "can't create get all trucks report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all trucks report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	if err = pdf.Output(ctx.Writer); err != nil {
		logger.ErrorLogger("CreateAllTrucksReport", "can't save get all trucks report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all trucks report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()
}
