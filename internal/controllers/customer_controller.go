package controllers

import (
	"bytes"
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

type CustomerController interface {
	CreateCustomer(ctx *gin.Context)
	GetCustomerById(ctx *gin.Context)
	GetAllCustomers(ctx *gin.Context)
	DeleteCustomer(ctx *gin.Context)
	UpdateCustomer(ctx *gin.Context)
	SearchCustomers(ctx *gin.Context)
	CreateAllCustomersReport(ctx *gin.Context)
}

type customerController struct {
	service         service.CustomerService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewCustomerController(service service.CustomerService, redisConnection *redis.Client, nats *nats.Conn) CustomerController {
	return &customerController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller *customerController) CreateCustomer(ctx *gin.Context) {

	var customer models.Customer
	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	err := validation.ValidateJSON(jsonData, customer)
	if err != nil {
		logger.ErrorLogger("CreateCustomer", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := ctx.BindJSON(&customer); err != nil {
		logger.ErrorLogger("CreateCustomer", "Can't bind customer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Customer json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := customer.IsValid(); err != nil {
		logger.ErrorLogger("CreateCustomer", "customer is not valid").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Customer is not valid", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_customer", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateCustomer", "Can't create customer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with crating Customer json", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Customer was created").Info(fmt.Sprintf("Customer was created with name %s.", customer.LegalName))

	ctx.JSON(http.StatusOK, "OK")

}

func (controller *customerController) GetCustomerById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	customer, err := controller.service.GetCustomerById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get customer by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return one customer", customer)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetCustomerById").Info(fmt.Sprintf("Customer was returned with id - %T", customer.ID))

}

func (controller *customerController) GetAllCustomers(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	customers, err := controller.service.GetAllCustomers(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get all customers", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all customers", customers)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllCustomers").Info("Return all customers")

}

func (controller *customerController) DeleteCustomer(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{
		"id":         ctx.Param("id"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	message := queue.Message{Method: "delete_customer", Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteCustomer", "Can't delete customer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't delete customer", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Customer was deleted").Info(fmt.Sprintf("Customer was deleted with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *customerController) UpdateCustomer(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var customer models.CustomerUpdateInput
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		logger.ErrorLogger("UpdateCustomer", "Can't bind CustomerUpdateInput json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding CustomerUpdateInput json", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_customer", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateCustomer", "Can't update customer").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't update customer", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	logger.InfoLogger("Customer was updated").Info(fmt.Sprintf("Customer was updated with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *customerController) SearchCustomers(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	customers, err := controller.service.SearchCustomers(searchText, offSet, limit)

	if err != nil {
		res := response.BuildErrorResponse("can't get all customers by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all customers by search", customers)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("SearchCustomers").Info(fmt.Sprintf("Return all customers by search. searchText - %s", searchText))

}

func (controller *customerController) CreateAllCustomersReport(ctx *gin.Context) {

	customers, err := controller.service.GetAllCustomersForReport()

	if err != nil {
		res := response.BuildErrorResponse("can't get all customers for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	customerModel := models.Customer{}
	header := reports.GetHeaderAsArray(customerModel)

	var data [][]string
	for _, v := range customers {
		customer := reports.GetValuesAsArray(v)
		data = append(data, customer)
	}

	pdf := reports.NewReport(1230, "Customer")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllCustomersReport", "Can't create get all customers report").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't create get all customers report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	err = pdf.Output(ctx.Writer)

	if err != nil {
		logger.ErrorLogger("CreateAllCustomersReport", "can't save get all customers report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all customers report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()
	logger.ErrorLogger("CreateAllCustomersReport", "Report was created").Info("Report was created")
}
