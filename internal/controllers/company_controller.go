package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/queue"
	"github.com/trucktrace/pkg/reports"
	"github.com/trucktrace/pkg/validation"
)

type CompanyController interface {
	CreateCompany(ctx *gin.Context)
	GetCompanyById(ctx *gin.Context)
	GetAllCompanies(ctx *gin.Context)
	DeleteCompany(ctx *gin.Context)
	UpdateCompany(ctx *gin.Context)
	SearchCompany(ctx *gin.Context)
	GetAllCompaniesForReport(ctx *gin.Context)
}

type companyController struct {
	service         service.CompanyService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewCompanyController(service service.CompanyService, redisConnection *redis.Client, nats *nats.Conn) CompanyController {
	return &companyController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

type Payload struct {
	Message []byte
}

func (company *companyController) CreateCompany(ctx *gin.Context) {

	form, _ := ctx.MultipartForm()

	if len(form.Value) == 0 {
		logger.ErrorLogger("CreateCompany", "Company body is empty").Error("Error - " + "Company body is empty")
		res := response.BuildErrorResponse("Company body is empty", "Company body is empty", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	jsonString := form.Value["json"][0]

	var companyModel models.Company

	err := validation.ValidateJSON([]byte(jsonString), companyModel)
	if err != nil {
		logger.ErrorLogger("CreateCompany", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := json.Unmarshal([]byte(jsonString), &companyModel); err != nil {
		logger.ErrorLogger("CreateCompany", "Can't bind Company json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Company json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["file"]

	for _, fileHeader := range files {

		fileName, err := company.service.SaveOrderLogo(fileHeader, ctx)

		if err != nil {
			logger.ErrorLogger("CreateCompany", "Can't save company logo").Error("Error - " + err.Error())
			res := response.BuildErrorResponse("can't save a company logo", err.Error(), fileHeader.Filename)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		companyModel.PictureName = fileName
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	jsonData, err := json.Marshal(companyModel)
	if err != nil {
		logger.ErrorLogger("CreateCompany", "Can't marshal Company").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with marshalling Company", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	message := queue.Message{Method: "create_company", Body: jsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(company.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateCompany", "Can't create company").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating company", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Company was created").Info(fmt.Sprintf("Comapny was created with name %s.", companyModel.LegalName))
	ctx.JSON(http.StatusOK, "OK")

}

func (company *companyController) DeleteCompany(ctx *gin.Context) {
	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{
		"id":         ctx.Param("id"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	companyImg, err := company.service.GetCompanyById(&helper.OneQueryParams{Id: ctx.Param("id")})

	if err != nil {
		logger.ErrorLogger("DeleteCompany", "Can't get company by id").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't get company by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if companyImg.PictureName != "" {
		err = os.Remove("upload/companyImg/" + companyImg.PictureName)
	}

	if err != nil {
		logger.ErrorLogger("DeleteCompany", "Can't delete company logo").Error("Error- " + err.Error())
		res := response.BuildErrorResponse("can't delete company logo", err.Error(), companyImg.PictureName)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	message := queue.Message{Method: "delete_company", Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(company.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteCompany", "Can't delete company").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't delete company", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Company was deleted").Info(fmt.Sprintf("Comapny was delete with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")
}

func (company *companyController) UpdateCompany(ctx *gin.Context) {

	form, _ := ctx.MultipartForm()
	var companyModel models.CompanyUpdateInput

	if len(form.Value) == 0 {
		logger.ErrorLogger("UpdateCompany", "Company body is empty").Error("Error - Order body is empty")
		res := response.BuildErrorResponse("Company body is empty", "Company body is empty", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	jsonString := form.Value["json"][0]

	if err := json.Unmarshal([]byte(jsonString), &companyModel); err != nil {
		logger.ErrorLogger("UpdateCompany", "Can't bind Company json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Company json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["file"]

	for _, fileHeader := range files {

		companyImg, _ := company.service.GetCompanyById(&helper.OneQueryParams{Id: ctx.Param("id")})

		if companyImg.PictureName != "" {
			if err := os.Remove("upload/companyImg/" + companyImg.PictureName); err != nil {
				logger.ErrorLogger("UpdateCompany", "Can't delete company logo").Error("Error - " + err.Error())
				res := response.BuildErrorResponse("can't delete company logo", err.Error(), fileHeader.Filename)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
				return
			}
		}

		fileName, err := company.service.SaveOrderLogo(fileHeader, ctx)

		if err != nil {
			logger.ErrorLogger("UpdateCompany", "Can't save company logo").Error("Error - " + err.Error())
			res := response.BuildErrorResponse("can't save company logo", err.Error(), fileHeader.Filename)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		companyModel.PictureName = fileName

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

	jsonData, err := json.Marshal(companyModel)

	if err != nil {
		logger.ErrorLogger("UpdateCompany", "Can't marshal Company").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with marshalling Company", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	message := queue.Message{Method: "update_company", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsUpdateRequest(company.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateCompany", "Can't update company").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't update company", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	logger.InfoLogger("Company was updated").Info(fmt.Sprintf("Comapny was updated with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")
}

func (company *companyController) GetCompanyById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	c, err := company.service.GetCompanyById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get company by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Return one company", c)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetCompanyById").Info("Return one company")
}

func (company *companyController) GetAllCompanies(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	companies, err := company.service.GetAllCompanies(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get all companies", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all companies", companies)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllCompanies").Info("Return all companies")

}

func (company *companyController) SearchCompany(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	companies, err := company.service.SearchCompany(searchText, offSet, limit)

	if err != nil {
		res := response.BuildErrorResponse("can't get all companies by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Return all customers by search", companies)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchCompany").Info(fmt.Sprintf("Return all company by search. searchText - %s", searchText))

}

func (company *companyController) GetAllCompaniesForReport(ctx *gin.Context) {

	companies, err := company.service.GetAllCompaniesForReport()

	if err != nil {
		res := response.BuildErrorResponse("can't get all companies for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	companyModel := models.Company{}
	header := reports.GetHeaderAsArray(companyModel)

	var data [][]string
	for _, v := range companies {
		company := reports.GetValuesAsArray(v)
		data = append(data, company)
	}

	pdf := reports.NewReport(980, "Company")
	pdf = reports.Header(pdf, header, 60)
	pdf = reports.Table(pdf, data, 60)

	if pdf.Err() {
		logger.ErrorLogger("GetAllCompaniesForReport", "Can't create get all companies report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all companies report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	err = pdf.Output(ctx.Writer)

	if err != nil {
		logger.ErrorLogger("GetAllCompaniesForReport", "Can't save get all companies report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't save get all companies report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()
	logger.InfoLogger("GetAllCompanies").Info("Report was created for all company")

}
