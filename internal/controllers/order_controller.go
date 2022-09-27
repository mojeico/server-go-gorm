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

type OrderController interface {
	CreateOrder(ctx *gin.Context)
	GetAllOrdersByCompanyId(ctx *gin.Context)
	GetAllOrders(ctx *gin.Context)
	GetOrderById(ctx *gin.Context)
	GetOrderByDriverId(ctx *gin.Context)
	GetOrderByTrailerId(ctx *gin.Context)
	GetOrderByTruckId(ctx *gin.Context)
	UpdateOrder(ctx *gin.Context)
	DeleteOrder(ctx *gin.Context)
	SearchOrders(ctx *gin.Context)

	MakeOrderCompleted(ctx *gin.Context)
	MakeOrderInvoiced(ctx *gin.Context)
	CreateAllOrdersReport(ctx *gin.Context)

	CreateOrderReportById(ctx *gin.Context)

	CreateOrderExtraPay(ctx *gin.Context)
	GetExtraPaysByOrderId(ctx *gin.Context)
}

type orderController struct {
	service         service.OrderService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewOrderController(service service.OrderService, redisConnection *redis.Client, nats *nats.Conn) OrderController {
	return &orderController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

//-----------------------+CREATE+-----------------------//
func (controller *orderController) CreateOrder(ctx *gin.Context) {

	var order models.Order
	if err := ctx.BindJSON(&order); err != nil {
		logger.ErrorLogger("CreateOrder", "Can't bind Order json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Order json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID
	order.CompanyId = companyId

	jsonData, _ := json.Marshal(order)

	err := validation.ValidateJSON(jsonData, order)
	if err != nil {
		logger.ErrorLogger("CreateOrder", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_order", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateOrder", "Problem with creating order").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating order", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Order was created").Info(fmt.Sprintf("Order was created with load nubmer %s.", order.LoadNumber))

	ctx.JSON(http.StatusOK, "OK")
}

//-----------------------+GET ALL ORDERS+-----------------------//
func (controller *orderController) GetAllOrders(ctx *gin.Context) {

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

	orders, err := controller.service.GetAllOrders(queryParams, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all orders", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("orders was gotten", orders)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllOrders").Info("Return all orders")

}

//-----------------------+GET ORDER BY COMPANY ID+-----------------------//
func (controller *orderController) GetAllOrdersByCompanyId(ctx *gin.Context) {

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

	orders, err := controller.service.GetAllOrdersByCompanyId(companyId, queryParams, orderBy, byte(orderDir))
	if err != nil {
		res := response.BuildErrorResponse("can't get all orders", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("orders was gotten", orders)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllOrdersByCompanyId").Info("Return all orders by company id")

}

//-----------------------+GET ORDER BY ID+-----------------------//
func (controller *orderController) GetOrderById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	order, err := controller.service.GetOrderById(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get order by id", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("order by id was gotten", order)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetOrderById").Info(fmt.Sprintf("Order was returned with id - %T", order.ID))

}

func (controller *orderController) GetOrderByDriverId(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("driverId"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	order, err := controller.service.GetOrderByDriverId(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get order by driverId", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("order by driverId was gotten", order)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetOrderByDriverId").Info("order by driverId was gotten")

}

func (controller *orderController) GetOrderByTrailerId(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("trailerId"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	order, err := controller.service.GetOrderByTrailerId(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get order by trailerId", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("order by trailerId was gotten", order)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetOrderByDriverId").Info("order by driverId was gotten")

}

func (controller *orderController) GetOrderByTruckId(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("truckId"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	order, err := controller.service.GetOrderByTruckId(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get order by truckId", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("order by truckId was gotten", order)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetOrderByTruckId").Info("order by truckId was gotten")

}

// -----------------------+UPDATE+-----------------------//
func (controller *orderController) UpdateOrder(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var orderUpdate models.OrderUpdateInput
	if err := ctx.BindJSON(&orderUpdate); err != nil {
		logger.ErrorLogger("UpdateOrder", "Can't bind OrderUpdateInput json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding OrderUpdateInput json", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_order", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateOrder", "Can't update Order").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating Order", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Order was updated").Info(fmt.Sprintf("Order was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

//-----------------------+DELETE+-----------------------//
func (controller *orderController) DeleteOrder(ctx *gin.Context) {

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

	message := queue.Message{Method: "delete_order", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteOrder", "Can't delete Order").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Order", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Order was deleted").Info(fmt.Sprintf("Order was deleted with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *orderController) SearchOrders(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	orders, err := controller.service.SearchOrders(searchText, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all orders by search", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("orders was gotten by search", orders)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("SearchOrders").Info(fmt.Sprintf("Return all orders by search. searchText - %s", searchText))

}

func (controller *orderController) MakeOrderCompleted(ctx *gin.Context) {

	var queries = map[string]string{

		"id": ctx.Param("id"),
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "order_to_completed", Body: nil, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("MakeOrderCompleted", "Can't make Order as completed").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with making Order as completed", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *orderController) MakeOrderInvoiced(ctx *gin.Context) {

	var queries = map[string]string{
		"id": ctx.Param("id"),
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "order_to_invoiced", Body: nil, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("MakeOrderInvoiced", "Can't make Order as invoicing").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with making Order as invoicing", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *orderController) CreateAllOrdersReport(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	orders, err := controller.service.GetAllOrdersForReport(companyId)

	if err != nil {

		res := response.BuildErrorResponse("can't get all orders for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	orderModel := models.Charges{}
	header := reports.GetHeaderAsArray(orderModel)

	var data [][]string
	for _, v := range orders {
		order := reports.GetValuesAsArray(v)
		data = append(data, order)
	}

	pdf := reports.NewReport(865, "Order")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllOrdersReport", "can't create get all orders report").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't create get all orders report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	if err = pdf.Output(ctx.Writer); err != nil {
		logger.ErrorLogger("CreateAllOrdersReport", "can't create get all orders report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all orders report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()
	logger.ErrorLogger("CreateAllOrdersReport", "Order report created").Info("Error - " + err.Error())

}

func (controller *orderController) CreateOrderReportById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	order, err := controller.service.GetOrderById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get order by id", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf := reports.CreateOrderReport(order)

	orderBytes, _ := pdf.Output()

	res := response.BuildResponse("Report was created", orderBytes)
	ctx.JSON(http.StatusOK, res)

	logger.ErrorLogger("CreateAllOrdersReport", "Order report created").Info("Error - " + err.Error())

}

func (controller *orderController) CreateOrderExtraPay(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID

	message := queue.Message{Method: "create_extrapay", Body: jsonData, UserID: userID}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateOrderExtraPay", "Can't create order extra pay").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating order extra pay", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *orderController) GetExtraPaysByOrderId(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("orderId"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	extraPays, err := controller.service.GetExtraPaysByOrderId(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get all extra pays", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("extra pays was gotten", extraPays)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetExtraPaysByOrderId").Info("extra pays was gotten")

}
