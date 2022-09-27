package test

/*
import (
	"bytes"
	"errors"
	"fmt"
	"github.com/trucktrace/internal/controllers"
	"github.com/trucktrace/internal/models"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	mock_service "github.com/trucktrace/internal/service/mocks"
)

var driverCreateOkBody string = `{
    "first_name": "testName",
    "last_name": "testLastName",
    "email": "test@gmail.com"
}`

var driverUpdateOkBody string = `{
				 "first_name": "newTestName",
    				"last_name": "newTestLastName"
        }`

func Test_driverController_CreateDriver(t *testing.T) {

	type mockBehavior func(s *mock_service.MockDriverService, driver models.Driver)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDriver          models.Driver
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Created",
			inputBody: driverCreateOkBody,
			inputDriver: models.Driver{
				FirstName: "testName",
				LastName:  "testLastName",
				Email:     "test@gmail.com",
			},
			mockBehavior: func(s *mock_service.MockDriverService, driver models.Driver) {
				s.EXPECT().CreateDriver(driver).Return(0, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"status":true,"message":"driver was created","errors":null,"data":0}`,
		},
		{
			name:      "Service or Database error Created",
			inputBody: driverCreateOkBody,
			inputDriver: models.Driver{
				FirstName: "testName",
				LastName:  "testLastName",
				Email:     "test@gmail.com",
			},
			mockBehavior: func(s *mock_service.MockDriverService, driver models.Driver) {
				s.EXPECT().CreateDriver(driver).Return(0, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't create a driver","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockDriverService(c)
			testCase.mockBehavior(service, testCase.inputDriver)

			controller := controllers.NewDriverController(service)

			r := gin.New()
			r.POST("/drivers", controller.CreateDriver)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/drivers",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}

}

func Test_driverController_DeleteDriver(t *testing.T) {
	type mockBehavior func(s *mock_service.MockDriverService, driverId, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputDriverId        int
		inputCompanyId       int
		inputStatus          string
		inputIsDeleted       string
		inputIsActive        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "OK Deleted",
			inputBody:      1,
			inputDriverId:  1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockDriverService, driverId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().DeleteDriver(driverId, companyId, status, isDeleted, isActive).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"driver was deleted","errors":null,"data":"ID: 1, CompanyID: 1"}`,
		},
		{
			name:           "Service or Database error Deleted",
			inputBody:      1,
			inputDriverId:  1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockDriverService, driverId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().DeleteDriver(driverId, companyId, status, isDeleted, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't delete driver","errors":"service or database error","data":1}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockDriverService(c)
			testCase.mockBehavior(
				service,
				testCase.inputDriverId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewDriverController(service)

			r := gin.New()
			r.DELETE("/drivers/:id", controller.DeleteDriver)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/drivers/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputDriverId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_driverController_GetAllDrivers(t *testing.T) {
	type mockBehavior func(s *mock_service.MockDriverService, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		inputStatus          string
		inputIsDeleted       string
		inputIsActive        string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK get all",

			mockBehavior: func(s *mock_service.MockDriverService, status, isDeleted, isActive string) {
				var driver = models.Driver{}
				driver.ID = 1
				s.EXPECT().GetAllDrivers(status, isDeleted, isActive).Return([]driver.Driver{driver}, nil)
			},
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"gotten all drivers succesfully","errors":null,"data":[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"first_name":"","last_name":"","address":"","city":"","state":"","zip":"","country":"","phone":"","email":"","gender":"","birth_day":"0001-01-01T00:00:00Z","status":"","is_active":false,"is_deleted":false,"checks":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"driver_id":0,"active":false,"freeze_payable":false,"gross_pay":false,"eligible":false,"allow_checks":false},"driver_attaches":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"driver_id":0,"company_id":0,"driver_group_id":0,"settlement_id":0,"employment_id":0,"safety_id":0}}]}`,
		},
		{
			name: "Service or Database error get all",
			mockBehavior: func(s *mock_service.MockDriverService, status, isDeleted, isActive string) {
				var driver = models.Driver{}
				driver.ID = 1
				s.EXPECT().GetAllDrivers(status, isDeleted, isActive).Return([]driver.Driver{driver}, errors.New("service or database error"))
			},
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't get all drivers","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockDriverService(c)
			testCase.mockBehavior(service,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewDriverController(service)

			r := gin.New()
			r.GET("/drivers/", controller.GetAllDrivers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/drivers/?status=%s&deleted=%s&active=%s",
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func Test_driverController_GetDriverById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockDriverService, driverId, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputDriverId        int
		inputCompanyId       int
		inputStatus          string
		inputIsDeleted       string
		inputIsActive        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK get by id",
			mockBehavior: func(s *mock_service.MockDriverService, driverId, companyId int, status, isDeleted, isActive string) {
				var driver = models.Driver{}
				driver.ID = 1
				s.EXPECT().GetDriverById(driverId, companyId, status, isDeleted, isActive).Return(driver, nil)
			},
			inputBody:            1,
			inputDriverId:        1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"gotten driver successfully","errors":null,"data":{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"first_name":"","last_name":"","address":"","city":"","state":"","zip":"","country":"","phone":"","email":"","gender":"","birth_day":"0001-01-01T00:00:00Z","status":"","is_active":false,"is_deleted":false,"checks":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"driver_id":0,"active":false,"freeze_payable":false,"gross_pay":false,"eligible":false,"allow_checks":false},"driver_attaches":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"driver_id":0,"company_id":0,"driver_group_id":0,"settlement_id":0,"employment_id":0,"safety_id":0}}}`,
		},
		{
			name: "Service or Database error get by id",
			mockBehavior: func(s *mock_service.MockDriverService, driverId, companyId int, status, isDeleted, isActive string) {
				var driver = models.Driver{}
				driver.ID = 1
				s.EXPECT().GetDriverById(driverId, companyId, status, isDeleted, isActive).Return(driver, errors.New("service or database error"))
			},
			inputBody:            1,
			inputDriverId:        1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't get a driver by id","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockDriverService(c)
			testCase.mockBehavior(service,
				testCase.inputDriverId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewDriverController(service)

			r := gin.New()
			r.GET("/drivers/:id", controller.GetDriverById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/drivers/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputDriverId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_driverController_UpdateDriver(t *testing.T) {
	type mockBehavior func(
		s *mock_service.MockDriverService,
		driverInput models.DriverUpdateInput,
		driverId, companyId int,
		status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputDriver          models.DriverUpdateInput
		inputDriverId        int
		inputCompanyId       int
		inputStatus          string
		inputIsDeleted       string
		inputIsActive        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Updated",
			inputBody: driverUpdateOkBody,
			inputDriver: models.DriverUpdateInput{
				FirstName: "newTestName",
				LastName:  "newTestLastName",
			},
			mockBehavior: func(s *mock_service.MockDriverService, driver models.DriverUpdateInput, driverId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().UpdateDriver(driverId, companyId, status, isDeleted, isActive, driver).Return(nil)
			},
			inputDriverId:        1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"driver was updated","errors":null,"data":"ID: 1, CompanyID: 1"}`,
		},
		{
			name:      "Service or Database error Updated",
			inputBody: driverUpdateOkBody,
			inputDriver: models.DriverUpdateInput{
				FirstName: "newTestName",
				LastName:  "newTestLastName",
			},
			mockBehavior: func(s *mock_service.MockDriverService, driver models.DriverUpdateInput, driverId int, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().UpdateDriver(driverId, companyId, status, isDeleted, isActive, driver).Return(errors.New("service or database error"))
			},
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't update driver","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockDriverService(c)
			testCase.mockBehavior(service,
				testCase.inputDriver,
				testCase.inputDriverId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewDriverController(service)

			r := gin.New()
			r.PATCH("/drivers/:id", controller.UpdateDriver)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/drivers/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputDriverId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}
*/
