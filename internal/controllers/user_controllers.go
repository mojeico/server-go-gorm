package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/trucktrace/pkg/logger"
	"github.com/trucktrace/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/response"
	"github.com/trucktrace/internal/service"
	"github.com/trucktrace/pkg/email"
	"github.com/trucktrace/pkg/helper"
	"github.com/trucktrace/pkg/queue"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetAllUsersByCompanyId(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	SearchUsers(ctx *gin.Context)
}

type userController struct {
	service         service.UserService
	redisConnection *redis.Client
	nats            *nats.Conn
}

func NewUserController(userService service.UserService, redisConnection *redis.Client, nats *nats.Conn) UserController {
	return &userController{
		service:         userService,
		redisConnection: redisConnection,
		nats:            nats,
	}
}

func (handler *userController) CreateUser(ctx *gin.Context) {

	form, _ := ctx.MultipartForm()

	if len(form.Value) == 0 {
		logger.ErrorLogger("CreateUser", "user body is empty").Error("Error - can't create a user with empty fields")
		res := response.BuildErrorResponse("user body is empty", "user body is empty", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	jsonString := form.Value["json"][0]

	var userModel models.User
	err := json.Unmarshal([]byte(jsonString), &userModel)
	if err != nil {
		logger.ErrorLogger("CreateUser", "can't binding User json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding User json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = userModel.IsValid()
	if err != nil {
		logger.ErrorLogger("CreateUser", "Can't create User").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with creating User", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["file"]

	var fileName string

	for _, fileHeader := range files {

		fileName, err = handler.service.SaveUserLogo(fileHeader, ctx)
		if err != nil {
			logger.ErrorLogger("CreateUser", "Can't save user logo").Error("Error - " + err.Error())
			res := response.BuildErrorResponse("can't save a user logo", err.Error(), response.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		userModel.PictureName = fileName
	}

	userModel.IpAddress = ctx.ClientIP()
	userModel.UserAgent = ctx.GetHeader("User-Agent")
	userModel.LastConnection = time.Now()

	userContext, _ := ctx.Get("currentUser")
	userID := userContext.(models.User).ID
	userEmail := userContext.(models.User).Email

	if err != nil {
		logger.ErrorLogger("CreateUser", "Token claims wrong type").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with token", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userModel.Status = "activated" // activated - confirm - deleted  -- login only for confirm
	userModel.Role = "USER"

	userJsonData, err := json.Marshal(&userModel)
	if err != nil {
		logger.ErrorLogger("CreateUser", "can't marshal json to user model").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't marshal json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = validation.ValidateJSON(userJsonData, userModel)
	if err != nil {
		logger.ErrorLogger("CreateUser", "can't validate json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("problem with validate json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	receiver := email.NewRequest([]string{userModel.Email}, "Registration successful")
	receiver.Send("./pkg/email/emailTemplate/createUser.html", map[string]string{"Username": userModel.Username, "Password": userModel.Password})

	message := queue.Message{Method: "create_user", Body: userJsonData, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsCreateRequest(handler.nats, message)

	if err != nil {
		logger.SystemLoggerError("CreateUser", "Can't create User").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't create User", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("User was created").Info(fmt.Sprintf("User was created with username %s.", userModel.Username))

	ctx.JSON(http.StatusOK, "OK")
}

func (handler *userController) DeleteUser(ctx *gin.Context) {

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

	companyImg, _ := handler.service.GetUserById(&helper.OneQueryParams{Id: ctx.Param("id")})

	var err error
	if companyImg.PictureName != "" {
		err = os.Remove("upload/userImg/" + companyImg.PictureName)
	}

	if err != nil {
		logger.ErrorLogger("DeleteUser", "Can't delete user logo").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't delete user logo", err.Error(), companyImg.PictureName)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	jsonData, _ := ioutil.ReadAll(ctx.Request.Body)

	message := queue.Message{Method: "delete_user", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsDeleteRequest(handler.nats, message)

	if err != nil {
		logger.SystemLoggerError("DeleteUser", "Can't delete User").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't delete User", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("User was deleted").Info(fmt.Sprintf("User was deleted with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")
}

func (handler *userController) UpdateUser(ctx *gin.Context) {

	form, _ := ctx.MultipartForm()

	if len(form.Value) == 0 {
		logger.ErrorLogger("UpdateUser", "User body is empty").Error("Error - empty user body")
		res := response.BuildErrorResponse("User body is empty", "user body is empty", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	jsonString := form.Value["json"][0]
	var userModel models.UpdateUserInput

	err := json.Unmarshal([]byte(jsonString), &userModel)
	if err != nil {
		logger.ErrorLogger("UpdateUser", "Can't bind user json").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with binding user json", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["file"]

	var fileName string

	for _, fileHeader := range files {

		companyImg, _ := handler.service.GetUserById(&helper.OneQueryParams{Id: ctx.Param("id")})

		if companyImg.PictureName != "" {
			if err = os.Remove("upload/userImg/" + companyImg.PictureName); err != nil {
				logger.ErrorLogger("UpdateUser", "Can't delete user logo").Error("Error - " + err.Error())
				res := response.BuildErrorResponse("can't delete user logo", err.Error(), fileHeader.Filename)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
				return
			}
		}

		fileName, err = handler.service.SaveUserLogo(fileHeader, ctx)

		if err != nil {
			logger.ErrorLogger("UpdateUser", "Can't save user logo").Error("Error - " + err.Error())
			res := response.BuildErrorResponse("can't save user logo", err.Error(), fileHeader.Filename)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		userModel.PictureName = fileName
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

	jsonData, err := json.Marshal(userModel)

	if err != nil {
		logger.ErrorLogger("UpdateUser", "Can't marshal user").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with marshalling user", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	message := queue.Message{Method: "update_user", Body: jsonData, Params: queries, UserID: userID, UserEmail: userEmail}

	serviceResponse, err := queue.RunNatsUpdateRequest(handler.nats, message)

	if err != nil {
		logger.SystemLoggerError("UpdateUser", "Can't update User").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Problem with updating User", serviceResponse.UserErrorMessage, response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	logger.InfoLogger("Customer was updated").Info(fmt.Sprintf("Customer was updated with id =  %s .", queries["id"]))

	ctx.JSON(http.StatusOK, "OK")

}

func (handler *userController) GetAllUsers(ctx *gin.Context) {
	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	queryParams := helper.NewGetAllQueryParams(
		ctx.Query("offset"),
		ctx.Query("limit"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
		ctx.Query("field"),
		ctx.Query("value"),
	)

	users, err := handler.service.GetAllUsers(queryParams, companyId)
	if err != nil {
		logger.ErrorLogger("GetAllUsers", "Can't get all users").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't get all users", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("gotten all users successfully", users)
	ctx.JSON(http.StatusOK, res)
	logger.InfoLogger("GetAllUsers").Info("User list was returned")
}

func (handler *userController) GetAllUsersByCompanyId(ctx *gin.Context) {
	companyId := ctx.Param("company_id")

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

	trucks, err := handler.service.GetAllUsersByCompanyId(companyId, queryParams, orderBy, byte(orderDir))
	if err != nil {
		logger.ErrorLogger("GetAllUsersByCompanyId", "can't get all users by user id").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("can't get all users", err.Error(), emptyObj)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponse("users by user id was gotten", trucks)
	ctx.JSON(http.StatusOK, res)
}

func (handler *userController) GetUserById(ctx *gin.Context) {
	queryParams := helper.NewOneQueryParams(
		ctx.Param("id"),
		ctx.Query("status"),
		ctx.Query("deleted"),
		ctx.Query("active"),
	)

	user, err := handler.service.GetUserById(queryParams)
	if err != nil {
		logger.ErrorLogger("GetUserById", "Can't get user by id").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't get user by id", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("User by id was gotten successfully", user)
	ctx.JSON(http.StatusOK, res)
}

func (handler *userController) SearchUsers(ctx *gin.Context) {

	offSet := ctx.Query("offset")
	limit := ctx.Query("limit")
	searchText := ctx.Query("text")

	userContext, _ := ctx.Get("currentUser")
	companyId := userContext.(models.User).CompanyID

	users, err := handler.service.SearchUsers(searchText, offSet, limit, companyId)

	if err != nil {
		logger.ErrorLogger("SearchUsers", "Can't get all users by search").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("Can't get all users by search", err.Error(), response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponse("gotten all users successfully by search", users)
	ctx.JSON(http.StatusOK, res)
}
