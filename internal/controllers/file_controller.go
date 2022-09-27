package controllers

import (
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
	"github.com/trucktrace/pkg/validation"
)

type FileController interface {
	UploadFile(ctx *gin.Context)
	GetFileById(ctx *gin.Context)
	GetLastFileByOwnerId(ctx *gin.Context)
	GetAllFilesByOwnerId(ctx *gin.Context)
	GetAllFiles(ctx *gin.Context)
	DeleteFile(ctx *gin.Context)
	SearchFiles(ctx *gin.Context)
}

type fileController struct {
	service         service.FileService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewFileController(service service.FileService, redisConnection *redis.Client, nats *nats.Conn) FileController {
	return &fileController{
		service:         service,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (file *fileController) UploadFile(ctx *gin.Context) {

	ownerType := fmt.Sprint(ctx.PostForm("owner"))
	ownerId, err := strconv.Atoi(fmt.Sprint(ctx.PostForm("owner_id")))

	if err != nil {
		logger.ErrorLogger("UploadFile", "can't convert owner id to number").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("owner id is not a number", err.Error(), "UploadFile")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	f, err := ctx.FormFile("file")

	if err != nil {
		logger.ErrorLogger("UploadFile", "Can't form file").Error("Error - " + err.Error())

		res := response.BuildErrorResponse("no file is received", err.Error(), f)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	expDate, err := strconv.Atoi(ctx.PostForm("expiration_date"))
	if err != nil {
		expDate = 0
	}

	expStatus := ctx.PostForm("expiration_status")
	comment := ctx.PostForm("comment")

	newFileName, ext, err := file.service.UploadFile(f.Filename)
	if err != nil {
		logger.ErrorLogger("UploadFile", "Can't upload file").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't upload a file", err.Error(), f.Filename)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	if err := ctx.SaveUploadedFile(f, "upload/"+newFileName); err != nil {
		logger.ErrorLogger("UploadFile", "Can't save file").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save a file", err.Error(), f.Filename)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	newFile := models.File{
		Name:             newFileName,
		Extension:        ext,
		Size:             f.Size,
		ExpirationDate:   int64(expDate),
		ExpirationStatus: expStatus,

		Comment: comment,

		OwnerType: ownerType,
		OwnerId:   ownerId,
	}
	newFileJsonData, _ := json.Marshal(newFile)

	err = validation.ValidateJSON(newFileJsonData, models.File{})
	if err != nil {
		logger.ErrorLogger("UploadFile", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	message := queue.Message{Method: "create_file", Body: newFileJsonData, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsCreateRequest(file.nats, message)

	if err != nil {
		logger.SystemLoggerError("UploadFile", "Can't save file in bd").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't save a file in bd ", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("File was uploaded").Info(fmt.Sprintf("File was uploaded with name =  %s .", newFile.Name))
	ctx.JSON(http.StatusOK, "OK")
}

func (file *fileController) GetFileById(ctx *gin.Context) {

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	f, err := file.service.GetFileById(queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get file by id", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("file by id was gotten", f)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetFileById").Info(fmt.Sprintf("Customer was returned with id - %T", f.ID))

}

func (file *fileController) GetLastFileByOwnerId(ctx *gin.Context) {

	ownerType := ctx.Query("owner")

	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	f, err := file.service.GetLastFileByOwnerId(ownerType, queryParams)

	if err != nil {
		res := response.BuildErrorResponse("can't get file by owner id", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("file by id was gotten", f)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetLastFileByOwnerId").Info("file by id was gotten")
}

func (file *fileController) GetAllFilesByOwnerId(ctx *gin.Context) {

	ownerType := ctx.Query("owner")

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	files, err := file.service.GetAllFilesByOwnerId(ownerType, queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get all files by owner", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("files by owner was gotten", files)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("GetAllFilesByOwnerId").Info("files by owner was gotten")
}

func (file *fileController) GetAllFiles(ctx *gin.Context) {

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	files, err := file.service.GetAllFiles(queryParams)
	if err != nil {
		res := response.BuildErrorResponse("can't get all files", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("files was gotten", files)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllFiles").Info("files was gotten")

}

func (file *fileController) DeleteFile(ctx *gin.Context) {

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	var queries = map[string]string{

		"id":         ctx.Param("id"),
		"owner_type": ctx.Query("owner"),
		"status":     ctx.Query("status"),
		"is_deleted": ctx.Query("deleted"),
		"is_active":  ctx.Query("active"),
	}

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)

	message := queue.Message{Method: "delete_file", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}
	serviceResponse, err := queue.RunNatsDeleteRequest(file.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteFile", "Can't delete file").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with deleting file ", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("File was deleted").Info(fmt.Sprintf("File was deleted with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")
}

func (file *fileController) SearchFiles(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	files, err := file.service.SearchFiles(searchText, offSet, limit)

	if err != nil {
		res := response.BuildErrorResponse("can't get all files by search", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("files was gotten by search", files)
	ctx.JSON(http.StatusOK, res)

	logger.InfoLogger("SearchFiles").Info(fmt.Sprintf("Return all files by search. searchText - %s", searchText))
}
