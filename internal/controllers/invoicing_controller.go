package controllers

import (
	"bytes"
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

type InvoicingController interface {
	CreateInvoicing(ctx *gin.Context)
	UpdateInvoicing(ctx *gin.Context)
	DeleteInvoicing(ctx *gin.Context)

	GetAllInvoices(ctx *gin.Context)
	GetInvoicingById(ctx *gin.Context)
	SearchAndFilterInvoices(ctx *gin.Context)
	CreateAllInvoicesReport(ctx *gin.Context)
	GetInvoiceInPDF(ctx *gin.Context)
}

type invoicingController struct {
	service         service.InvoicingService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewInvoicingController(service service.InvoicingService, redisConnection *redis.Client, nats *nats.Conn) InvoicingController {
	return &invoicingController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (controller *invoicingController) CreateInvoicing(ctx *gin.Context) {
	var invoicing models.Invoicing

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	err := validation.ValidateJSON(jsonData, invoicing)
	if err != nil {
		logger.ErrorLogger("CreateInvoicing", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := ctx.BindJSON(&invoicing); err != nil {
		logger.ErrorLogger("CreateInvoicing", "Can't bind Invoicing json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Invoicing json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_invoicing", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateInvoicing", "Can't create Invoicing").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating Invoicing", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Invoicing was created").Info("Invoicing was created")

	ctx.JSON(http.StatusOK, "OK")
}

func (controller *invoicingController) UpdateInvoicing(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var invoiceUpdate models.InvoicingUpdateInput
	if err := ctx.BindJSON(&invoiceUpdate); err != nil {
		logger.ErrorLogger("UpdateInvoicing", "Can't bind InvoicingUpdateInput").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding InvoicingUpdateInput json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{
		"id": ctx.Param("id"),
	}

	message := queue.Message{Method: "update_invoicing", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateInvoicing", "Can't update Invoicing").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't update Invoicing", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Invoicing was updated").Info(fmt.Sprintf("Invoicing was updated with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *invoicingController) DeleteInvoicing(ctx *gin.Context) {
	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{
		"id": ctx.Param("id"),
	}

	message := queue.Message{Method: "delete_invoicing", Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(controller.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteInvoicing", "Can't delete Invoicing").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Invoicing", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Invoicing was deleted").Info(fmt.Sprintf("Invoicing was deleted with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (controller *invoicingController) GetAllInvoices(ctx *gin.Context) {

	queryParams := helper.NewPaginationParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("deleted"),
	)

	invoices, err := controller.service.GetAllInvoices(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get all Invoices", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all invoices", invoices)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllInvoices").Info("Get all group by company id successfully")

}

func (controller *invoicingController) GetInvoicingById(ctx *gin.Context) {

	id := ctx.Param("id")

	invoices, err := controller.service.GetInvoicingById(id)

	if err != nil {
		res := response.BuildErrorResponse("can't get invoicing by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return one invoicing", invoices)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetInvoicingById").Info(fmt.Sprintf("Invoicing was returned with id - %T", invoices.ID))

}

func (controller *invoicingController) SearchAndFilterInvoices(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	query := helper.NewInvoicingFilter(
		strconv.Itoa(companyId),
		ctx.Query("deliveryFrom"),
		ctx.Query("deliveryTo"),
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("deleted"),
	)

	settlements, err := controller.service.SearchAndFilterInvoices(query)

	if err != nil {
		res := response.BuildErrorResponse("can't get all invoices by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all invoices  by search", settlements)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("SearchAndFilterInvoices").Info("Return all invoices  by search")

}

func (controller *invoicingController) CreateAllInvoicesReport(ctx *gin.Context) {

	invoices, err := controller.service.GetAllInvoicesForReport()

	if err != nil {
		res := response.BuildErrorResponse("can't get all invoices for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	invoiceModel := models.Invoicing{}
	header := reports.GetHeaderAsArray(invoiceModel)

	var data [][]string
	for _, v := range invoices {
		invoice := reports.GetValuesAsArray(v)
		data = append(data, invoice)
	}

	pdf := reports.NewReport(870, "Invoicing")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllInvoicesReport", "can't create get all invoices report").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't create get all invoices report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	if err = pdf.Output(ctx.Writer); err != nil {
		logger.ErrorLogger("CreateAllInvoicesReport", "can't save get all invoices report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all invoices report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()

	logger.ErrorLogger("CreateAllInvoicesReport", "Invoices report was crated").Info("Error - " + err.Error())

}

func (controller *invoicingController) GetInvoiceInPDF(ctx *gin.Context) {
	id := ctx.Param("id")

	invoice, err := controller.service.GetInvoicingById(id)

	if err != nil {
		logger.ErrorLogger("GetInvoiceInPDF", "can't create pdf for invoice for this id").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("can't create pdf for invoice for this id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf := reports.CreateInvoiceReport(invoice)

	invoicingBytes, err := pdf.Output()

	if err != nil {
		logger.ErrorLogger("GetInvoiceInPDF", "can't generate invoice pdf file").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't generate invoice pdf file", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Invoice pdf was created", invoicingBytes)
	ctx.JSON(http.StatusOK, res)
	logger.ErrorLogger("GetInvoiceInPDF", "Invoice pdf was created").Info("Error - " + err.Error())

}
