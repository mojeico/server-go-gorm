package controllers

import (
	"github.com/trucktrace/pkg/logger"
	"net/http"
	"strconv"

	"github.com/trucktrace/internal/privileges"

	"github.com/trucktrace/internal/models"
	"github.com/trucktrace/internal/repository"
	"github.com/trucktrace/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/trucktrace/internal/response"
)

type PrivilegeController interface {
	CheckPrivilege(ctx *gin.Context)
	GetAllPrivileges(ctx *gin.Context)
}

type privilegeController struct {
	userRepo  repository.UserRepository
	groupRepo repository.GroupRepository
}

func NewPrivilegeController(userRepo repository.UserRepository, groupRepo repository.GroupRepository) PrivilegeController {
	return &privilegeController{
		userRepo:  userRepo,
		groupRepo: groupRepo,
	}
}

func (handler *privilegeController) CheckPrivilege(ctx *gin.Context) {

	value, isCookie := ctx.Request.Header["Set-Cookie"]

	if !isCookie {
		logger.ErrorLogger("CheckPrivilege", "Cookie is empty").Error("Header Set-Cookie is empty")

		res := response.BuildErrorResponse("User is not authorized", "Cookie is empty", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	claims, err := token.ParseToken(value[0])

	if err != nil {
		logger.ErrorLogger("CheckPrivilege", "Problem with parse token").Error("Error - " + err.Error())
		res := response.BuildErrorResponse("User is not authorized", "Cookie is empty", response.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}

	user, _ := handler.userRepo.GetUserByIdForMiddleware(strconv.Itoa(claims.UserId))

	var groups []models.Groups

	for _, value := range user.Groups {
		group, _ := handler.groupRepo.GetGroupByIdForMiddleware(strconv.Itoa(int(value)))
		groups = append(groups, group)
	}

	var privileges []string

	for _, group := range groups {
		for _, val := range group.Priveleges {
			privileges = append(privileges, val)
		}
	}

	res := response.BuildErrorResponse("Get all user privileges", "Get all user privileges", privileges)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)

	logger.InfoLogger("CheckPrivilege").Info("Return all user privileges list")

}

func (handler *privilegeController) GetAllPrivileges(ctx *gin.Context) {

	privelege := privileges.Privilege{

		SearchGroups:            "search_groups",
		GetAllGroupsByCompanyId: "get_all_groups_by_company_id",
		GetAllGroups:            "get_all_groups",
		GetGroupById:            "get_group_by_id",
		UpdateGroup:             "update_group",
		CrateGroup:              "create_group",
		DeleteGroup:             "delete_group",

		SearchUsers:            "search_users",
		GetAllUsersByCompanyId: "get_all_users_by_company_id",
		GetAllUsers:            "get_all_users",
		GetUserById:            "get_user_by_id",
		UpdateUser:             "update_user",
		CrateUser:              "create_user",
		DeleteUser:             "delete_user",

		SearchDrivers:            "search_drivers",
		GetAllDriversByCompanyId: "get_all_drivers_by_company_id",
		GetAllDrivers:            "get_all_drivers",
		GetDriverById:            "get_driver_by_id",
		UpdateDriver:             "update_driver",
		CrateDriver:              "create_driver",
		DeleteDriver:             "delete_driver",

		SearchTrailers:            "search_trailers",
		GetAllTrailersByCompanyId: "get_all_trailers_by_company_id",
		GetAllTrailers:            "get_all_trailers",
		GetTrailerById:            "get_trailer_by_id",
		UpdateTrailer:             "update_trailer",
		CrateTrailer:              "create_trailer",
		DeleteTrailer:             "delete_trailer",

		SearchSettlements: "search_settlements",
		GetAllSettlement:  "get_all_settlements",
		GetSettlementById: "get_settlement_by_id",
		CreateSettlement:  "create_settlement",
		DeleteSettlement:  "delete_settlement",
		UpdateSettlement:  "update_settlement",

		SearchTrucks:            "search_trucks",
		GetAllTrucksByCompanyId: "get_all_trucks_by_company_id",
		GetAllTrucks:            "get_all_trucks",
		GetTruckById:            "get_truck_by_id",
		UpdateTruck:             "update_truck",
		CrateTruck:              "create_truck",
		DeleteTruck:             "delete_truck",

		SearchCustomers: "search_customers",
		GetAllCustomers: "get_all_customers",
		GetCustomerById: "get_customer_by_id",
		UpdateCustomer:  "update_customer",
		CrateCustomer:   "create_customer",
		DeleteCustomer:  "delete_customer",

		SearchSafeties:            "search_safeties",
		GetAllSafetiesByCompanyId: "get_all_safeties_by_company_id",
		GetAllSafeties:            "get_all_safeties",
		GetSafetyById:             "get_safety_by_id",
		UpdateSafety:              "update_safety",
		CrateSafety:               "create_safety",
		DeleteSafety:              "delete_safety",

		SearchOrders:            "search_orders",
		GetAllOrdersByCompanyId: "get_all_orders_by_company_id",
		GetAllOrders:            "get_all_orders",
		GetOrderById:            "get_order_by_id",
		UpdateOrder:             "update_order",
		CrateOrder:              "create_order",
		DeleteOrder:             "delete_order",

		GetAllCompanies: "get_all_companies",
		GetCompanyById:  "get_company_by_id",
		UpdateCompany:   "update_company",
		CrateCompany:    "create_company",
		DeleteCompany:   "delete_company",

		CheckPrivilege:   "get_all_privileges_by_user",
		GetAllPrivileges: "get_all_privileges",

		UploadFile: "upload_file",
	}

	res := response.BuildErrorResponse("Get all privileges", "Get all privileges", privelege)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
}
