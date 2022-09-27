package privileges

type Privilege struct {
	SearchGroups            string `json:"search_groups"`
	GetAllGroupsByCompanyId string `json:"get_all_groups_by_company_id"`
	GetAllGroups            string `json:"get_all_groups"`
	GetGroupById            string `json:"get_group_by_id"`
	UpdateGroup             string `json:"update_group"`
	CrateGroup              string `json:"crate_group"`
	DeleteGroup             string `json:"delete_group"`

	SearchUsers            string `json:"search_users"`
	GetAllUsersByCompanyId string `json:"get_all_users_by_company_id"`
	GetAllUsers            string `json:"get_all_users"`
	GetUserById            string `json:"get_user_by_id"`
	UpdateUser             string `json:"update_user"`
	CrateUser              string `json:"crate_user"`
	DeleteUser             string `json:"delete_user"`

	SearchDrivers            string `json:"search_drivers"`
	GetAllDriversByCompanyId string `json:"get_all_drivers_by_company_id"`
	GetAllDrivers            string `json:"get_all_drivers"`
	GetDriverById            string `json:"get_driver_by_id"`
	UpdateDriver             string `json:"update_driver"`
	CrateDriver              string `json:"crate_driver"`
	DeleteDriver             string `json:"delete_driver"`

	SearchTrailers            string `json:"search_trailers"`
	GetAllTrailersByCompanyId string `json:"get_all_trailers_by_company_id"`
	GetAllTrailers            string `json:"get_all_trailers"`
	GetTrailerById            string `json:"get_trailer_by_id"`
	UpdateTrailer             string `json:"update_trailer"`
	CrateTrailer              string `json:"crate_trailer"`
	DeleteTrailer             string `json:"delete_trailer"`

	SearchSettlements string `json:"search_settlements"`
	GetAllSettlement  string `json:"get_all_settlement"`
	GetSettlementById string `json:"get_settlement_by_id"`
	CreateSettlement  string `json:"create_settlement"`
	DeleteSettlement  string `json:"delete_settlement"`
	UpdateSettlement  string `json:"update_settlement"`

	SearchTrucks            string `json:"search_trucks"`
	GetAllTrucksByCompanyId string `json:"get_all_trucks_by_company_id"`
	GetAllTrucks            string `json:"get_all_trucks"`
	GetTruckById            string `json:"get_truck_by_id"`
	UpdateTruck             string `json:"update_truck"`
	CrateTruck              string `json:"crate_truck"`
	DeleteTruck             string `json:"delete_truck"`

	SearchCustomers string `json:"search_customers"`
	GetAllCustomers string `json:"get_all_customers"`
	GetCustomerById string `json:"get_customer_by_id"`
	UpdateCustomer  string `json:"update_customer"`
	CrateCustomer   string `json:"crate_customer"`
	DeleteCustomer  string `json:"delete_customer"`

	SearchSafeties            string `json:"search_safeties"`
	GetAllSafetiesByCompanyId string `json:"get_all_safeties_by_company_id"`
	GetAllSafeties            string `json:"get_all_safeties"`
	GetSafetyById             string `json:"get_safety_by_id"`
	UpdateSafety              string `json:"update_safety"`
	CrateSafety               string `json:"crate_safety"`
	DeleteSafety              string `json:"delete_safety"`

	SearchOrders            string `json:"search_orders"`
	GetAllOrdersByCompanyId string `json:"get_all_orders_by_company_id"`
	GetAllOrders            string `json:"get_all_orders"`
	GetOrderById            string `json:"get_order_by_id"`
	UpdateOrder             string `json:"update_order"`
	CrateOrder              string `json:"crate_order"`
	DeleteOrder             string `json:"delete_order"`

	SearchCompanies string `json:"search_companies"`
	GetAllCompanies string `json:"get_all_companies"`
	GetCompanyById  string `json:"get_company_by_id"`
	UpdateCompany   string `json:"update_company"`
	CrateCompany    string `json:"crate_company"`
	DeleteCompany   string `json:"delete_company"`

	CheckPrivilege   string `json:"get_all_privileges_by_user"`
	GetAllPrivileges string `json:"get_all_privileges"`

	UploadFile string `json:"upload_file"`
}

// type Priveleges interface {
// 	CreateUser()
// 	UpdateUser()
// 	GetAllUsers()
// 	GetUserById()
// }

// type PrivelegesList = interface{}{
// 	"CreateUser":  true,
// 	"DeleteUser":  true,
// 	"GetAllUsers": true,
// }

// type PrivelegesStruct struct {
// 	Priveleges PrivelegesList
// }

// type CompanyDirector struct {
// 	AllUsers
// 	AllPrivelegs
// 	AllGroupsCreated
// 	AllModels
// }

// type Groups struct {
// 	Name string
// 	Director.id string
// 	PrivelegesList map[string]bool
// 	Company.id
// 	user.ids
// }

// type Priveleges struct {
// 	Id int
// 	PrivName string
// 	State bool
// }
