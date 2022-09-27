package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nats-io/nats.go"

	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/queue"
	"github.com/trucktrace/pkg/reports"
	"github.com/trucktrace/pkg/validation"

	"github.com/gin-gonic/gin"
)

type GroupController interface {
	CreateGroup(ctx *gin.Context)
	DeleteGroup(ctx *gin.Context)
	UpdateGroup(ctx *gin.Context)
	GetGroupByID(ctx *gin.Context)
	GetAllGroups(ctx *gin.Context)
	GetAllGroupsByCompanyId(ctx *gin.Context)
	SearchGroups(ctx *gin.Context)
	CreateAllGroupsReport(ctx *gin.Context)
}

type groupController struct {
	service         service.GroupService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewGroupController(groupService service.GroupService, redisConnection *redis.Client, nats *nats.Conn) GroupController {
	return &groupController{
		service:         groupService,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (handler *groupController) GetGroupByID(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	group, err := handler.service.GetGroupById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("Can not find group by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("Get group by id successfully", group)
	ctx.AbortWithStatusJSON(http.StatusOK, res)
	logger.InfoLogger("GetGroupByID").Info("Get group by id successfully")

}

func (handler *groupController) CreateGroup(ctx *gin.Context) {

	var group models.Groups

	if err := ctx.BindJSON(&group); err != nil {
		logger.ErrorLogger("CreateGroup", "Problem with binding Groups json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding Groups json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID
	group.CompanyId = strconv.Itoa(companyId)

	jsonData, _ := json.Marshal(group)
	err := validation.ValidateJSON(jsonData, group)
	if err != nil {
		logger.ErrorLogger("CreateGroup", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	message := queue.Message{Method: "create_group", Body: jsonData, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsCreateRequest(handler.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateGroup", "Problem with creating Group").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating Group", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Group was created").Info(fmt.Sprintf("Group was created with name %s.", group.Name))

	ctx.JSON(http.StatusOK, "OK")
}

func (handler *groupController) DeleteGroup(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{
		"id":         ctx.Param("id"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	message := queue.Message{Method: "delete_group", Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsDeleteRequest(handler.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteGroup", "Can't delete group").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Group", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Group was deleted").Info(fmt.Sprintf("Group was deleted with id =  %s .", queries["id"]))
	ctx.JSON(http.StatusOK, "OK")

}

func (handler *groupController) UpdateGroup(ctx *gin.Context) {

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonData))

	var group models.GroupUpdateInput
	if err := ctx.BindJSON(&group); err != nil {
		logger.ErrorLogger("UpdateGroup", "Can't bind GroupUpdateInput json").Error("Error  - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding GroupUpdateInput json", err.Error(), response.EmptyObj{})
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

	message := queue.Message{Method: "update_group", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsUpdateRequest(handler.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteGroup", "Can't delete group").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting Group", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Group was updated").Info(fmt.Sprintf("Group was updated with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")

}

func (handler *groupController) GetAllGroups(ctx *gin.Context) {

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

	groups, err := handler.service.GetAllGroups(queryParams, companyId)

	if err != nil {
		res := response.BuildErrorResponse("Can't get all group", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Get all groups successfully", groups)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllGroups").Info("Get all groups successfully")

}

func (handler *groupController) GetAllGroupsByCompanyId(ctx *gin.Context) {

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

	orderBy := ctx.Query("orderBy")
	orderDir, _ := strconv.Atoi(ctx.Query("orderDir"))

	groups, err := handler.service.GetAllGroupsByCompanyId(strconv.Itoa(companyId), queryParams, orderBy, byte(orderDir))

	if err != nil {
		res := response.BuildErrorResponse("Can't get all group by company id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Get all group by company id successfully", groups)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllGroupsByCompanyId").Info("Get all group by company id successfully")

}

func (handler *groupController) SearchGroups(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	groups, err := handler.service.SearchGroups(searchText, offSet, limit, companyId)

	if err != nil {
		res := response.BuildErrorResponse("Can't get all groups by search ", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("Get all group by search successfully", groups)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("SearchGroups").Info("Get all group by search successfully")

}

func (handler *groupController) CreateAllGroupsReport(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	groups, err := handler.service.GetAllGroupsForReport(companyId)

	if err != nil {
		res := response.BuildErrorResponse("can't get all groups for report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	groupModel := models.Groups{}
	header := reports.GetHeaderAsArray(groupModel)

	var data [][]string
	for _, v := range groups {
		group := reports.GetValuesAsArray(v)
		data = append(data, group)
	}

	pdf := reports.NewReport(530, "Group")
	pdf = reports.Header(pdf, header, 50)
	pdf = reports.Table(pdf, data, 50)

	if pdf.Err() {
		logger.ErrorLogger("CreateAllGroupsReport", "can't create get all groups report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't create get all groups report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	err = pdf.Output(ctx.Writer)

	if err != nil {
		logger.ErrorLogger("CreateAllGroupsReport", "can't save get all groups report").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save get all groups report", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	pdf.Close()

	logger.ErrorLogger("CreateAllGroupsReport", "Report was created").Error("Error - " + err.Error())

}
