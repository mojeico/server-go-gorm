package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trucktrace/internal/controllers"
)

type AppApi struct {
	driverController         controllers.DriverController
	userController           controllers.UserController
	groupController          controllers.GroupController
	settlementController     controllers.SettlementController
	truckController          controllers.TruckController
	customerController       controllers.CustomerController
	trailerController        controllers.TrailerController
	safetyController         controllers.SafetyController
	orderController          controllers.OrderController
	privilegeController      controllers.PrivilegeController
	companyController        controllers.CompanyController
	chargesController        controllers.ChargesController
	invoicingController      controllers.InvoicingController
	fileController           controllers.FileController
	trailerCommentController controllers.TrailerCommentController
}

func NewAppApi(
	driverController controllers.DriverController,
	userController controllers.UserController,
	groupController controllers.GroupController,
	settlementController controllers.SettlementController,
	truckController controllers.TruckController,
	customerController controllers.CustomerController,
	trailerController controllers.TrailerController,
	safetyController controllers.SafetyController,
	orderController controllers.OrderController,
	privilegeController controllers.PrivilegeController,
	companyController controllers.CompanyController,
	chargesController controllers.ChargesController,
	invoicingController controllers.InvoicingController,
	fileController controllers.FileController,
	trailerCommentController controllers.TrailerCommentController,

) *AppApi {
	return &AppApi{
		driverController:         driverController,
		userController:           userController,
		groupController:          groupController,
		settlementController:     settlementController,
		truckController:          truckController,
		customerController:       customerController,
		trailerController:        trailerController,
		safetyController:         safetyController,
		orderController:          orderController,
		privilegeController:      privilegeController,
		companyController:        companyController,
		chargesController:        chargesController,
		invoicingController:      invoicingController,
		fileController:           fileController,
		trailerCommentController: trailerCommentController,
	}
}

// CreateDriver godoc
// @Summary Create Driver
// @Description Send Driver data to create
// @Tags driver
// @Accept  json
// @Produce  json
// @Param driver body models.Driver true "Driver"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /driver/ [post]
func (api *AppApi) CreateDriver(ctx *gin.Context) {
	api.driverController.CreateDriver(ctx)
}

// GetDriverById godoc
// @Summary Get Driver by id
// @Description Get all Driver data by id
// @Tags driver
// @Accept  json
// @Produce  json
// @Param  id path string true "Driver id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drivers/driver/:id [get]
func (api *AppApi) GetDriverById(ctx *gin.Context) {
	api.driverController.GetDriverById(ctx)

}

// GetAllDrivers godoc
// @Summary Get all Drivers
// @Description Get data about all Drivers
// @Tags driver
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /driver [get]
func (api *AppApi) GetAllDrivers(ctx *gin.Context) {
	api.driverController.GetAllDrivers(ctx)
}

// GetAllDriversByCompanyId godoc
// @Summary Get all Drivers by company id
// @Description Get data about all Drivers
// @Tags driver
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drivers/:companyId [get]
func (api *AppApi) GetAllDriversByCompanyId(ctx *gin.Context) {
	api.driverController.GetAllDriversByCompanyId(ctx)
}

// DeleteDriver godoc
// @Summary Create Driver
// @Description Send Driver data to create
// @Tags driver
// @Accept  json
// @Produce  json
// @Param  id path string true "Driver id"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /driver/:id [delete]
func (api *AppApi) DeleteDriver(ctx *gin.Context) {
	api.driverController.DeleteDriver(ctx)
}

// UpdateDriver godoc
// @Summary Update Driver
// @Description Update Driver data
// @Tags driver
// @Accept  json
// @Produce  json
// @Param  id path string true "Driver id"
// @Param driver body models.DriverUpdateInput true "Driver Update Input"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /driver/:id [patch]
func (api *AppApi) UpdateDriver(ctx *gin.Context) {
	api.driverController.UpdateDriver(ctx)
}

func (api *AppApi) CreateAllDriversReport(ctx *gin.Context) {
	api.driverController.CreateAllDriversReport(ctx)
}

// SearchDrivers godoc
// @Summary Get all Drivers by search
// @Description Get all Drivers by search
// @Tags driver
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /drivers/search [get]
func (api *AppApi) SearchDrivers(ctx *gin.Context) {
	api.driverController.SearchDrivers(ctx)
}

// GetGroupByID godoc
// @Summary Get Group by id
// @Description Get Group data by id
// @Tags group
// @Accept  json
// @Produce  json
// @Param  groupId path string true "Group id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /groups/group/:id [get]
func (api *AppApi) GetGroupByID(ctx *gin.Context) {
	api.groupController.GetGroupByID(ctx)
}

// CreateGroup godoc
// @Summary Create Group
// @Description Send Group data to create
// @Tags group
// @Accept  json
// @Produce  json
// @Param group body models.Groups true "Groups"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /group [post]
func (api *AppApi) CreateGroup(ctx *gin.Context) {
	api.groupController.CreateGroup(ctx)
}

// DeleteGroup godoc
// @Summary Delete Group
// @Description Sent Group id to delete it
// @Tags group
// @Accept  json
// @Produce  json
// @Param  groupId path string true "Group id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /group/:groupId [delete]
func (api *AppApi) DeleteGroup(ctx *gin.Context) {
	api.groupController.DeleteGroup(ctx)
}

// UpdateGroup godoc
// @Summary Update Group
// @Description Sent Group data to upgrade it
// @Tags group
// @Accept  json
// @Produce  json
// @Param group body models.GroupUpdateInput true "GroupUpdateInput"
// @Success 201 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /group/:groupId [post]
func (api *AppApi) UpdateGroup(ctx *gin.Context) {
	api.groupController.UpdateGroup(ctx)
}

// GetAllGroups godoc
// @Summary Get all Groups
// @Description  Get all Groups data
// @Tags group
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Router /group [get]
func (api *AppApi) GetAllGroups(ctx *gin.Context) {
	api.groupController.GetAllGroups(ctx)

}

// GetAllGroupsByCompanyId godoc
// @Summary Get all Groups by company
// @Description Get all Groups data by company
// @Tags group
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Param  companyId path string false "Company id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /group/:companyId [get]
func (api *AppApi) GetAllGroupsByCompanyId(ctx *gin.Context) {
	api.groupController.GetAllGroupsByCompanyId(ctx)
}

func (api *AppApi) CreateAllGroupsReport(ctx *gin.Context) {
	api.groupController.CreateAllGroupsReport(ctx)
}

// SearchGroups godoc
// @Summary Get all Groups by search
// @Description Get all Groups by search
// @Tags group
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /group/search [get]
func (api *AppApi) SearchGroups(ctx *gin.Context) {
	api.groupController.SearchGroups(ctx)
}

// CreateSettlement godoc
// @Summary Create Settlement
// @Description Create settlement based on orderId
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param orderId path string true "order orderId"
// @Success 201 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /settlement/:orderId [post]
func (api *AppApi) CreateSettlement(ctx *gin.Context) {
	api.settlementController.CreateSettlement(ctx)
}

// DeleteSettlement godoc
// @Summary Delete Settlement
// @Description Delete Settlement by id
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param  id path string true "Settlement id"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /settlement/:id [delete]
func (api *AppApi) DeleteSettlement(ctx *gin.Context) {
	api.settlementController.DeleteSettlement(ctx)

}

// UpdateSettlement godoc
// @Summary Update Settlement
// @Description Sent Settlement data to create
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param settlement body models.SettlementUpdateInput true "SettlementUpdateInput"
// @Param  id path string true "Settlement id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /settlement/:id [post]
func (api *AppApi) UpdateSettlement(ctx *gin.Context) {
	api.settlementController.UpdateSettlement(ctx)
}

// GetAllSettlement godoc
// @Summary Get all Settlements
// @Description Get all Settlements data
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Param settlement body models.Settlement true "Settlement"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /settlement [get]
func (api *AppApi) GetAllSettlement(ctx *gin.Context) {
	api.settlementController.GetAllSettlement(ctx)
}

// GetSettlementById godoc
// @Summary Get Settlement by id
// @Description Bet Settlement data by id
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Param  id path string true "Settlement id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /settlement/:id [get]
func (api *AppApi) GetSettlementById(ctx *gin.Context) {
	api.settlementController.GetSettlementById(ctx)
}

//TODO
func (api *AppApi) CreateAllSettlementsReport(ctx *gin.Context) {
	api.settlementController.CreateAllSettlementsReport(ctx)
}

// SearchSettlements godoc
// @Summary Get all Settlements by search
// @Description Get all Settlements by search
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param  text query string false "object search by text"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /settlement/search [get]
func (api *AppApi) SearchSettlements(ctx *gin.Context) {
	api.settlementController.SearchSettlements(ctx)
}

// CreateUser godoc
// @Summary Create User
// @Description Send User data to create
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body models.User true "User"
// @Success 201 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user [post]
func (api *AppApi) CreateUser(ctx *gin.Context) {
	api.userController.CreateUser(ctx)
}

// DeleteUser godoc
// @Summary Delete User
// @Description Delete User by id
// @Tags user
// @Accept  json
// @Produce  json
// @Param  id path string true "User id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/:id [delete]
func (api *AppApi) DeleteUser(ctx *gin.Context) {
	api.userController.DeleteUser(ctx)
}

// UpdateUser godoc
// @Summary Update User
// @Description Send User data to create
// @Tags user
// @Accept  json
// @Produce  json
// @Param UpdateUserInput body models.UpdateUserInput true "UpdateUserInput"
// @Param  id path string true "User id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/:id [post]
func (api *AppApi) UpdateUser(ctx *gin.Context) {
	api.userController.UpdateUser(ctx)
}

// GetAllUsers godoc
// @Summary Get all Users
// @Description Get all User data
// @Tags user
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Param User body models.User true "User"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /users [get]
func (api *AppApi) GetAllUsers(ctx *gin.Context) {
	api.userController.GetAllUsers(ctx)
}

// GetAllUsersByCompanyId godoc
// @Summary Get all Users by Company id
// @Description Get all User data
// @Tags user
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Param User body models.User true "User"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /users/:companyId [get]
func (api *AppApi) GetAllUsersByCompanyId(ctx *gin.Context) {
	api.userController.GetAllUsersByCompanyId(ctx)
}

// GetUserById godoc
// @Summary Get User by id
// @Description Get User data by id
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/user/:id [get]
func (api *AppApi) GetUserById(ctx *gin.Context) {
	api.userController.GetUserById(ctx)
}

// SearchUsers godoc
// @Summary Get all Users by search
// @Description Get all Users by search
// @Tags user
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/search [get]
func (api *AppApi) SearchUsers(ctx *gin.Context) {
	api.userController.SearchUsers(ctx)
}

// CreateTruck godoc
// @Summary Create Truck
// @Description Send Truck data to create
// @Tags trucks
// @Accept  json
// @Produce  json
// @Param truck body models.Truck true "Truck"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trucks/ [post]
func (api *AppApi) CreateTruck(ctx *gin.Context) {
	api.truckController.CreateTruck(ctx)
}

// GetAllTrucks godoc
// @Summary Get all Trucks
// @Description Get all Trucks data
// @Tags truck
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Param id path string true "Company Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trucks/ [get]
func (api *AppApi) GetAllTrucks(ctx *gin.Context) {
	api.truckController.GetAllTrucks(ctx)
}

// GetAllTrucksByCompanyId godoc
// @Summary Get all Truck by Company id
// @Description Get all User data
// @Tags truck
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Param User body models.Truck true "Truck"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /trucks/:companyId [get]
func (api *AppApi) GetAllTrucksByCompanyId(ctx *gin.Context) {
	api.truckController.GetAllTrucksByCompanyId(ctx)
}

// GetTruckById godoc
// @Summary Get Truck by id
// @Description Get Truck data by id
// @Tags truck
// @Accept  json
// @Produce  json
// @Param id path string true "truck id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trucks/truck/:id [get]
func (api *AppApi) GetTruckById(ctx *gin.Context) {
	api.truckController.GetTruckById(ctx)
}

// UpdateTruck godoc
// @Summary Update Truck
// @Description Update Truck Info
// @Tags truck
// @Accept  json
// @Produce  json
// @Param  id path string true "Truck Id"
// @Param truck body models.TruckUpdateInput true "TruckUpdateInput"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trucks/:truckId [put]
func (api *AppApi) UpdateTruck(ctx *gin.Context) {
	api.truckController.UpdateTruck(ctx)
}

//TODO
func (api *AppApi) CreateAllTrucksReport(ctx *gin.Context) {
	api.truckController.CreateAllTrucksReport(ctx)
}

// DeleteTruck godoc
// @Summary Delete Truck
// @Description Delete Truck Info
// @Tags truck
// @Accept  json
// @Produce  json
// @Param  id path string true "Truck Id"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trucks/:truckId [delete]
func (api *AppApi) DeleteTruck(ctx *gin.Context) {
	api.truckController.DeleteTruck(ctx)
}

// SearchTrucks godoc
// @Summary Get all Trucks by search
// @Description Get all Trucks by search
// @Tags truck
// @Accept  json
// @Produce  json
// @Param  text query string false "object search by text"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trucks/search [get]
func (api *AppApi) SearchTrucks(ctx *gin.Context) {
	api.truckController.SearchTrucks(ctx)
}

// CreateCustomer godoc
// @Summary Create Customer
// @Description Create Customer Info
// @Tags customer
// @Accept  json
// @Produce  json
// @Param customer body models.Customer true "Customer"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customer/ [post]
func (api *AppApi) CreateCustomer(ctx *gin.Context) {
	api.customerController.CreateCustomer(ctx)
}

// GetCustomerById godoc
// @Summary Get Customer By Id
// @Description Get Customer By Id saved in PostgreSQL
// @Tags truck
// @Accept  json
// @Produce  json
// @Param  id path string true "Customer Id"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customer/:id [get]
func (api *AppApi) GetCustomerById(ctx *gin.Context) {
	api.customerController.GetCustomerById(ctx)
}

// GetAllCustomers godoc
// @Summary Update Customer
// @Description Update Customer Info
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path string true "Customer Id"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customer [get]
func (api *AppApi) GetAllCustomers(ctx *gin.Context) {
	api.customerController.GetAllCustomers(ctx)
}

// DeleteCustomer godoc
// @Summary Delete Customer
// @Description Delete Customer Info
// @Tags customer
// @Accept  json
// @Produce  json
// @Param  id path string true "Customer Id"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customer/:customerId [delete]
func (api *AppApi) DeleteCustomer(ctx *gin.Context) {
	api.customerController.DeleteCustomer(ctx)
}

// UpdateCustomer godoc
// @Summary Update Customer
// @Description Update Customer Info
// @Tags customer
// @Accept  json
// @Produce  json
// @Param  id path string true "Customer Id"
// @Param truck body models.CustomerUpdateInput true "CustomerUpdateInput"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers/:customerId [put]
func (api *AppApi) UpdateCustomer(ctx *gin.Context) {
	api.customerController.UpdateCustomer(ctx)
}

func (api *AppApi) CreateAllCustomersReport(ctx *gin.Context) {
	api.customerController.CreateAllCustomersReport(ctx)
}

// SearchCustomers godoc
// @Summary Get all Customers by search
// @Description Get all Customers by search
// @Tags customer
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /customers/search [get]
func (api *AppApi) SearchCustomers(ctx *gin.Context) {
	api.customerController.SearchCustomers(ctx)
}

// GetAllTrailers godoc
// @Summary Update Trailer
// @Description Update Trailer Info
// @Tags trailer
// @Accept  json
// @Produce  json
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailer/ [get]
func (api *AppApi) GetAllTrailers(ctx *gin.Context) {
	api.trailerController.GetAllTrailers(ctx)
}

// GetAllTrailersByCompanyId godoc
// @Summary Get Trailer By Id
// @Description Get Trailer By Id saved in PostgreSQL
// @Tags trailer
// @Accept  json
// @Produce  json
// @Param companyId path string true "Company Id"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailer/company/:companyId [get]
func (api *AppApi) GetAllTrailersByCompanyId(ctx *gin.Context) {
	api.trailerController.GetAllTrailersByCompanyId(ctx)
}

// GetTrailerByID godoc
// @Summary Get Trailer By Id
// @Description Get Trailer By Id saved in PostgreSQL
// @Tags trailer
// @Accept  json
// @Produce  json
// @Param  id path string true "Trailer Id"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailers/trailer/:trailerId [get]
func (api *AppApi) GetTrailerByID(ctx *gin.Context) {
	api.trailerController.GetTrailerByID(ctx)
}

// CreateTrailer godoc
// @Summary Create Trailer
// @Description Send Trailer data to create
// @Tags trailers
// @Accept  json
// @Produce  json
// @Param trailer body models.Trailer true "Trailer"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailer/ [post]
func (api *AppApi) CreateTrailer(ctx *gin.Context) {
	api.trailerController.CreateTrailer(ctx)
}

// UpdateTrailer godoc
// @Summary Update Trailer
// @Description Update Trailer Info
// @Tags trailer
// @Accept  json
// @Produce  json
// @Param  id path string true "Trailer Id"
// @Param truck body models.TrailerInputUpdate true "TrailerUpdateInput"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailer/:trailerId [put]
func (api *AppApi) UpdateTrailer(ctx *gin.Context) {
	api.trailerController.UpdateTrailer(ctx)
}

// DeleteTrailer godoc
// @Summary Update Trailer
// @Description Update Trailer Info
// @Tags trailer
// @Accept  json
// @Produce  json
// @Param  id path string true "Trailer Id"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailer/:trailerId [delete]
func (api *AppApi) DeleteTrailer(ctx *gin.Context) {
	api.trailerController.DeleteTrailer(ctx)
}

func (api *AppApi) CreateAllTrailersReport(ctx *gin.Context) {
	api.trailerController.CreateAllTrailersReport(ctx)
}

// SearchTrailers godoc
// @Summary Get all Trailers by search
// @Description Get all Trailers by search
// @Tags trailer
// @Accept  json
// @Produce  json
// @Param  text query string false "object search by text"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailer/search [get]
func (api *AppApi) SearchTrailers(ctx *gin.Context) {
	api.trailerController.SearchTrailers(ctx)
}

// DeleteSafety godoc
// @Summary Update Safety
// @Description Delete Safety Info
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "Safety Id"
// @Param  companyId query string false "company id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safety/:id [delete]
func (api *AppApi) DeleteSafety(ctx *gin.Context) {
	api.safetyController.DeleteSafety(ctx)
}

// CreateSafety godoc
// @Summary Create Safety
// @Description Send Safety data to create
// @Tags safety
// @Accept  json
// @Produce  json
// @Param safety body models.Safety true "Safety"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safety [post]
func (api *AppApi) CreateSafety(ctx *gin.Context) {
	api.safetyController.CreateSafety(ctx)
}

// GetSafetyById godoc
// @Summary Get Safety By Id
// @Description Get Safety By Id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "Safety Id"
// @Param  companyId query string false "company id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safeties/:id [get]
func (api *AppApi) GetSafetyById(ctx *gin.Context) {
	api.safetyController.GetSafetyById(ctx)
}

// GetAllSafeties godoc
// @Summary Get all Safeties
// @Description Get All Safeties
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safety [get]
func (api *AppApi) GetAllSafeties(ctx *gin.Context) {
	api.safetyController.GetAllSafeties(ctx)
}

func (api *AppApi) GetAllSafetiesByType(ctx *gin.Context) {
	api.safetyController.GetAllSafetiesByType(ctx)
}

// GetAllSafetiesByCompanyId godoc
// @Summary Get all Safeties by company id
// @Description Get All Safeties by company id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  companyId path string true "Company Id"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safety/:companyId [get]
func (api *AppApi) GetAllSafetiesByCompanyId(ctx *gin.Context) {
	api.safetyController.GetAllSafetiesByCompanyId(ctx)
}

// UpdateSafety godoc
// @Summary Get all Safeties by company id
// @Description Get All Safeties by company id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "Safety Id"
// @Param  companyId query string true "company id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safety/:id [get]
func (api *AppApi) UpdateSafety(ctx *gin.Context) {
	api.safetyController.UpdateSafety(ctx)
}

func (api *AppApi) CreateAllSafetyReport(ctx *gin.Context) {
	api.safetyController.CreateAllSafetyReport(ctx)
}

// SearchSafety godoc
// @Summary Get all Safeties by search
// @Description Get all Safeties by search
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  text query string false "object search by text"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /safety/search [get]
func (api *AppApi) SearchSafety(ctx *gin.Context) {
	api.safetyController.SearchSafety(ctx)
}

// CreateOrder godoc
// @Summary Create order
// @Description Send Order data to create
// @Tags order
// @Accept  json
// @Produce  json
// @Param order body models.Order true "Order"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /order [post]
func (api *AppApi) CreateOrder(ctx *gin.Context) {
	api.orderController.CreateOrder(ctx)
}

//TODO
func (api *AppApi) CreateOrderExtraPay(ctx *gin.Context) {
	api.orderController.CreateOrderExtraPay(ctx)
}

// GetAllOrdersByCompanyId godoc
// @Summary Get all Orders by Company id
// @Description Get all orders data
// @Tags order
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit offset"
// @Param Order body models.Order true "Order"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /orders/:companyId [get]
func (api *AppApi) GetAllOrdersByCompanyId(ctx *gin.Context) {
	api.orderController.GetAllOrdersByCompanyId(ctx)
}

// GetAllOrders godoc
// @Summary Get all Orders
// @Description Get All Orders
// @Tags order
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders [get]
func (api *AppApi) GetAllOrders(ctx *gin.Context) {
	api.orderController.GetAllOrders(ctx)
}

// GetOrderById godoc
// @Summary Get Order By Id
// @Description Get Order By Id
// @Tags order
// @Accept  json
// @Produce  json
// @Param  id path string true "order Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/order/:id [get]
func (api *AppApi) GetOrderById(ctx *gin.Context) {
	api.orderController.GetOrderById(ctx)
}

//TODO
func (api *AppApi) GetExtraPaysByOrderId(ctx *gin.Context) {
	api.orderController.GetExtraPaysByOrderId(ctx)
}

// GetOrderByDriverId godoc
// @Summary Get Order By driver Id
// @Description Get Order By driver Id
// @Tags order
// @Accept  json
// @Produce  json
// @Param  driverId path string true "order driverId"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/order/driver/:driverId [get]
func (api *AppApi) GetOrderByDriverId(ctx *gin.Context) {
	api.orderController.GetOrderByDriverId(ctx)
}

// GetOrderByTrailerId godoc
// @Summary Get Order By trailer Id
// @Description Get Order By trailer Id
// @Tags order
// @Accept  json
// @Produce  json
// @Param  trailerId path string true "order trailerId"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/order/trailer/:trailerId [get]
func (api *AppApi) GetOrderByTrailerId(ctx *gin.Context) {
	api.orderController.GetOrderByTrailerId(ctx)
}

// GetOrderByTruckId godoc
// @Summary Get Order By truck Id
// @Description Get Order By truck Id
// @Tags order
// @Accept  json
// @Produce  json
// @Param  truckId path string true "order truckId"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/order/truck/:truckId [get]
func (api *AppApi) GetOrderByTruckId(ctx *gin.Context) {
	api.orderController.GetOrderByTruckId(ctx)
}

// UpdateOrder godoc
// @Summary Update Order by Id
// @Description Update Order by Id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "Order Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/:id [patch]
func (api *AppApi) UpdateOrder(ctx *gin.Context) {
	api.orderController.UpdateOrder(ctx)
}

// DeleteOrder godoc
// @Summary Delete Order by Id
// @Description Delete Order by Id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "Safety Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/:id [delete]
func (api *AppApi) DeleteOrder(ctx *gin.Context) {
	api.orderController.DeleteOrder(ctx)
}

// SearchOrders godoc
// @Summary Get all Orders by search
// @Description Get all Orders by search
// @Tags order
// @Accept  json
// @Produce  json
// @Param  text query string false "object search by text"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/search [get]
func (api *AppApi) SearchOrders(ctx *gin.Context) {
	api.orderController.SearchOrders(ctx)
}

// MakeOrderCompleted godoc
// @Summary Mark order completed by id
// @Description Mark order completed by id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "order id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/completed/:id [post]
func (api *AppApi) MakeOrderCompleted(ctx *gin.Context) {
	api.orderController.MakeOrderCompleted(ctx)
}

// MakeOrderInvoiced godoc
// @Summary Make order invoiced by id
// @Description Make order invoiced by id
// @Tags safety
// @Accept  json
// @Produce  json
// @Param  id path string true "order id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/invoiced/:id [post]
func (api *AppApi) MakeOrderInvoiced(ctx *gin.Context) {
	api.orderController.MakeOrderInvoiced(ctx)
}

// CreateAllOrdersReport godoc
// @Summary Create all orders report
// @Description Create all orders report
// @Tags order
// @Accept  json
// @Produce  json
// @Param  text query string false "object search by text"
// @Param  offset query string false "offset search by offset "
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders/report [get]
func (api *AppApi) CreateAllOrdersReport(ctx *gin.Context) {
	api.orderController.CreateAllOrdersReport(ctx)
}

func (api *AppApi) CheckPrivilege(ctx *gin.Context) {
	api.privilegeController.CheckPrivilege(ctx)
}

func (api *AppApi) GetAllPrivileges(ctx *gin.Context) {
	api.privilegeController.GetAllPrivileges(ctx)
}

// CreateCompany godoc
// @Summary Create company
// @Description Create company
// @Tags company
// @Accept  json
// @Produce  json
// @Param order body models.Company true "Company"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company [post]
func (api *AppApi) CreateCompany(ctx *gin.Context) {
	api.companyController.CreateCompany(ctx)
}

// SearchCompany godoc
// @Summary Get all companies by search
// @Description Get all companies by search
// @Tags company
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/search [get]
func (api *AppApi) SearchCompany(ctx *gin.Context) {
	api.companyController.SearchCompany(ctx)
}

// GetCompanyById godoc
// @Summary Get company By Id
// @Description Get company By Id
// @Tags company
// @Accept  json
// @Produce  json
// @Param  id path string true "Company Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/:id [get]
func (api *AppApi) GetCompanyById(ctx *gin.Context) {
	api.companyController.GetCompanyById(ctx)
}

// GetAllCompanies godoc
// @Summary Get all companies
// @Description Get All companies
// @Tags company
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company [get]
func (api *AppApi) GetAllCompanies(ctx *gin.Context) {
	api.companyController.GetAllCompanies(ctx)
}

// DeleteCompany godoc
// @Summary Delete Company by Id
// @Description Delete Company by Id
// @Tags company
// @Accept  json
// @Produce  json
// @Param  id path string true "Company Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/:id [delete]
func (api *AppApi) DeleteCompany(ctx *gin.Context) {
	api.companyController.DeleteCompany(ctx)
}

func (api *AppApi) GetAllCompaniesForReport(ctx *gin.Context) {
	api.companyController.GetAllCompaniesForReport(ctx)
}

// UpdateCompany godoc
// @Summary Update Company by Id
// @Description Update Company by Id
// @Tags company
// @Accept  json
// @Produce  json
// @Param  id path string true "Company Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /company/:id [patch]
func (api *AppApi) UpdateCompany(ctx *gin.Context) {
	api.companyController.UpdateCompany(ctx)
}

// CreateCharges godoc
// @Summary Create charges
// @Description Create charges
// @Tags charges
// @Accept  json
// @Produce  json
// @Param order body models.Charges true "Charges"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges [post]
func (api *AppApi) CreateCharges(ctx *gin.Context) {
	api.chargesController.CreateCharges(ctx)
}

// UpdateCharges godoc
// @Summary Update charges
// @Description Update charges
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  id path string true "charge Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges/:id [patch]
func (api *AppApi) UpdateCharges(ctx *gin.Context) {
	api.chargesController.UpdateCharges(ctx)
}

// DeleteCharges godoc
// @Summary Delete charges
// @Description Delete charges
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  id path string true "change Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges/:id [delete]
func (api *AppApi) DeleteCharges(ctx *gin.Context) {
	api.chargesController.DeleteCharges(ctx)
}

// GetChargeById godoc
// @Summary Get charge By Id
// @Description Get charge By Id
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  id path string true "charge Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges/charge/:id [get]
func (api *AppApi) GetChargeById(ctx *gin.Context) {
	api.chargesController.GetChargeById(ctx)
}

// GetAllChargesByOrderId godoc
// @Summary Get all charges by OrderId
// @Description Get all charges by OrderId
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  orderId path string true "order orderId"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges/order/:orderId [get]
func (api *AppApi) GetAllChargesByOrderId(ctx *gin.Context) {
	api.chargesController.GetAllChargesByOrderId(ctx)
}

// GetAllChargesBySettlementId godoc
// @Summary Get all charges by settlement Id
// @Description Get all charges by settlement Id
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  settlementId path string true "settlement settlementId"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /settlement/:settlementId [get]
func (api *AppApi) GetAllChargesBySettlementId(ctx *gin.Context) {
	api.chargesController.GetAllChargesBySettlementId(ctx)
}

// GetAllCharges godoc
// @Summary Get all charges
// @Description Get all charges
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit offset"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges [get]
func (api *AppApi) GetAllCharges(ctx *gin.Context) {
	api.chargesController.GetAllCharges(ctx)
}

// SearchCharge godoc
// @Summary Get all charges by search
// @Description Get all charges by search
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges/search [get]
func (api *AppApi) SearchCharge(ctx *gin.Context) {
	api.chargesController.SearchCharge(ctx)
}

// CreateAllChargersReport godoc
// @Summary Create all chargers report
// @Description Create all chargers report
// @Tags charges
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /charges/report [get]
func (api *AppApi) CreateAllChargersReport(ctx *gin.Context) {
	api.chargesController.CreateAllChargersReport(ctx)
}

//TODO
func (api *AppApi) CreateInvoicing(ctx *gin.Context) {
	api.invoicingController.CreateInvoicing(ctx)
}

// UpdateInvoicing godoc
// @Summary update invoice
// @Description update invoice
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  id path string true "invoice Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing/:id [patch]
func (api *AppApi) UpdateInvoicing(ctx *gin.Context) {
	api.invoicingController.UpdateInvoicing(ctx)
}

// DeleteInvoicing godoc
// @Summary delete invoice
// @Description delete invoice
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  id path string true "invoice Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing/:id [delete]
func (api *AppApi) DeleteInvoicing(ctx *gin.Context) {
	api.invoicingController.DeleteInvoicing(ctx)
}

// GetInvoicingById godoc
// @Summary Get invoice by id
// @Description Get invoice by id
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  id path string true "invoice Id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing/:id [get]
func (api *AppApi) GetInvoicingById(ctx *gin.Context) {
	api.invoicingController.GetInvoicingById(ctx)
}

// GetAllInvoices godoc
// @Summary Get all invoices
// @Description Get all invoices
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing [get]
func (api *AppApi) GetAllInvoices(ctx *gin.Context) {
	api.invoicingController.GetAllInvoices(ctx)
}

// SearchAndFilterInvoices godoc
// @Summary Search and filter invoices
// @Description Search and filter invoices
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  text query string false "text search by text"
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing/search [get]
func (api *AppApi) SearchAndFilterInvoices(ctx *gin.Context) {
	api.invoicingController.SearchAndFilterInvoices(ctx)
}

// CreateAllInvoicesReport godoc
// @Summary Create all invoices report
// @Description Create all invoices report
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  offset query string false "offset search by offset"
// @Param  limit query string false "limit search by limit"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing/report [get]
func (api *AppApi) CreateAllInvoicesReport(ctx *gin.Context) {
	api.invoicingController.CreateAllInvoicesReport(ctx)
}

// GetSettlementReportForDriverPDF godoc
// @Summary Search and filter invoices
// @Description Search and filter invoices if order is completed
// @Tags settlement
// @Accept  json
// @Produce  json
// @Param  id path string false "settlement id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /settlement/report/pdf/:id [get]
func (api *AppApi) GetSettlementReportForDriverPDF(ctx *gin.Context) {
	api.settlementController.GetSettlementReportForDriverPDF(ctx)
}

// GetInvoiceInPDF godoc
// @Summary Get invoice in PDF
// @Description Get invoice in PDF if order is completed
// @Tags invoice
// @Accept  json
// @Produce  json
// @Param  id path string false "settlement id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /invoicing/report/pdf/:id [get]
func (api *AppApi) GetInvoiceInPDF(ctx *gin.Context) {
	api.invoicingController.GetInvoiceInPDF(ctx)
}

// CreateTrailerComment godoc
// @Summary Create company
// @Description Create company
// @Tags trailer_comment
// @Accept  json
// @Produce  json
// @Param trailer body models.TrailerComment true "TrailerComment"
// @Param trailerId path string false "trailer id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailers/comment/:trailerId [post]
func (api *AppApi) CreateTrailerComment(ctx *gin.Context) {
	api.trailerCommentController.CreateTrailerComment(ctx)
}

// GetTrailerCommentsByTrailerID godoc
// @Summary Get trailer comments by trailer iD
// @Description Get trailer comments by trailer iD
// @Tags trailer_comment
// @Accept  json
// @Produce  json
// @Param id path string false "trailer id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /trailers/comment/:id [get]
func (api *AppApi) GetTrailerCommentsByTrailerID(ctx *gin.Context) {
	api.trailerCommentController.GetTrailerCommentsByTrailerID(ctx)
}
