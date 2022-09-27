package test

/*
import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/trucktrace/internal/models"
	mock_service "github.com/trucktrace/internal/service/mocks"
)

var settlementCreateOkBody string = `{
    "settlement_date":1622739708073,
    "invoicing_company":"test_company",
    "driver": 123,
    "total_miles":1,
    "empty_miles":1,
    "loaded_miles":1,
    "date_submitted":1622732738073,
    "deduction":1123.23,
    "reimbursement":4.3434,
    "earning":234.43,
    "total":24324.3
}`

var settlementUpdateOkBody string = `{
            "invoicing_company": "test_company_update"
        }`

func Test_settlementController_CreateSettlement(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSettlementService, settlement models.Settlement)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSettlement      models.Settlement
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Created",
			inputBody: settlementCreateOkBody,
			inputSettlement: models.Settlement{
				SettlementDate:   1622739708073,
				InvoicingCompany: "test_company",
				Driver:           123,
				TotalMiles:       1,
				EmptyMiles:       1,
				LoadedMiles:      1,
				DateSubmitted:    1622732738073,
				Deductions:       1123.23,
				Reimbursement:    4.3434,
				Earning:          234.43,
				Total:            24324.3,
			},
			mockBehavior: func(s *mock_service.MockSettlementService, settlement models.Settlement) {
				s.EXPECT().CreateSettlement(settlement).Return(0, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"status":true,"message":"Settlement created succesfully, id: 0","errors":null,"data":0}`,
		},
		{
			name:      "Service or Database error created",
			inputBody: settlementCreateOkBody,
			inputSettlement: models.Settlement{
				SettlementDate:   1622739708073,
				InvoicingCompany: "test_company",
				Driver:           123,
				TotalMiles:       1,
				EmptyMiles:       1,
				LoadedMiles:      1,
				DateSubmitted:    1622732738073,
				Deductions:       1123.23,
				Reimbursement:    4.3434,
				Earning:          234.43,
				Total:            24324.3,
			},
			mockBehavior: func(s *mock_service.MockSettlementService, settlement models.Settlement) {
				s.EXPECT().CreateSettlement(settlement).Return(0, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't create settlement ","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSettlementService(c)
			testCase.mockBehavior(service, testCase.inputSettlement)

			controllers := NewSettlementController(service)

			r := gin.New()
			r.POST("/settlement", controllers.CreateSettlement)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/settlement",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_settlementController_DeleteSettlement(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSettlementService, id int, status, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputSettlement      int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
	}{
		{
			name:            "OK Deleted",
			inputBody:       1,
			inputSettlement: 1,
			mockBehavior: func(s *mock_service.MockSettlementService, id int, status, isActive string) {
				s.EXPECT().DeleteSettlement(id, status, isActive).Return(nil)
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"Settlement deleted successfully, id: 1","errors":null,"data":1}`,
		},
		{
			name:            "Service or Database error Deleted",
			inputBody:       1,
			inputSettlement: 1,
			status:          "started",
			isActive:        "true",
			deleted:         "false",

			mockBehavior: func(s *mock_service.MockSettlementService, id int, status, isActive string) {
				s.EXPECT().DeleteSettlement(id, status, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can not delete settlement","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSettlementService(c)
			testCase.mockBehavior(service, testCase.inputSettlement, testCase.status, testCase.isActive)

			controllers := NewSettlementController(service)

			r := gin.New()
			r.DELETE("/settlement/:id", controllers.DeleteSettlement)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/settlement/%b?status=%s&active=%s&deleted=%s", testCase.inputBody, testCase.status, testCase.isActive, testCase.deleted),
				bytes.NewBufferString(fmt.Sprint(testCase.inputBody)))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_settlementController_GetAllSettlement(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSettlementService, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
	}{
		{
			name: "OK get all",

			mockBehavior: func(s *mock_service.MockSettlementService, status, isDeleted, isActive string) {
				s.EXPECT().GetAllSettlements(status, isDeleted, isActive).Return([]models.Settlement{{Driver: 1}}, nil)
			},
			expectedStatusCode:   200,
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedResponseBody: `{"status":true,"message":"Return all settlements","errors":null,"data":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"settlement_date":0,"invoicing_company":"","driver":1,"total_miles":0,"empty_miles":0,"loaded_miles":0,"date_submitted":0,"deduction":0,"reimbursement":0,"earning":0,"total":0,"status":"","deleted":false,"is_active":false}]}`,
		},
		{
			name: "Service or Database error get all",
			mockBehavior: func(s *mock_service.MockSettlementService, status, isDeleted, isActive string) {
				s.EXPECT().GetAllSettlements(status, isDeleted, isActive).Return([]models.Settlement{{Driver: 1}}, errors.New("service or database error"))
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't get all settlements","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSettlementService(c)
			testCase.mockBehavior(service, testCase.status, testCase.deleted, testCase.isActive)

			controllers := NewSettlementController(service)

			r := gin.New()
			r.GET("/settlement", controllers.GetAllSettlement)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/settlement?status=%s&active=%s&deleted=%s", testCase.status, testCase.isActive, testCase.deleted),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_settlementController_GetSettlementById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSettlementService, id int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSettlement      models.SettlementUpdateInput
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
	}{
		{
			name: "OK get by id",
			mockBehavior: func(s *mock_service.MockSettlementService, id int, status, isDeleted, isActive string) {
				s.EXPECT().GetSettlementById(id, status, isDeleted, isActive).Return(models.Settlement{InvoicingCompany: "Invoicing Company"}, nil)
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			inputId:              1,
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"Return one settlement","errors":null,"data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"settlement_date":0,"invoicing_company":"Invoicing Company","driver":0,"total_miles":0,"empty_miles":0,"loaded_miles":0,"date_submitted":0,"deduction":0,"reimbursement":0,"earning":0,"total":0,"status":"","deleted":false,"is_active":false}}`,
		},
		{
			name: "Service or Database error get by id",
			mockBehavior: func(s *mock_service.MockSettlementService, id int, status, isDeleted, isActive string) {
				s.EXPECT().GetSettlementById(id, status, isDeleted, isActive).Return(models.Settlement{}, errors.New("service or database error"))
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			inputId:              1,
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":false,"message":"can't get settlement by id","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSettlementService(c)
			testCase.mockBehavior(service, testCase.inputId, testCase.status, testCase.deleted, testCase.isActive)

			controllers := NewSettlementController(service)

			r := gin.New()
			r.GET("/settlement/:id", controllers.GetSettlementById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/settlement/%b?status=%s&active=%s&deleted=%s", testCase.inputId, testCase.status, testCase.isActive, testCase.deleted),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_settlementController_UpdateSettlement(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSettlementService, settlement models.SettlementUpdateInput, id int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSettlement      models.SettlementUpdateInput
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
	}{
		{
			name:      "OK Updated",
			inputBody: settlementUpdateOkBody,
			inputSettlement: models.SettlementUpdateInput{
				InvoicingCompany: "test_company_update",
			},
			status:   "started",
			isActive: "true",
			deleted:  "false",
			mockBehavior: func(s *mock_service.MockSettlementService, settlement models.SettlementUpdateInput, id int, status, isDeleted, isActive string) {
				s.EXPECT().UpdateSettlement(id, settlement, status, isDeleted, isActive).Return(nil)
			},
			inputId:              1,
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"settlement was updated","errors":null,"data":1}`,
		},
		{
			name:      "Service or Database error Updated",
			inputBody: settlementUpdateOkBody,
			inputSettlement: models.SettlementUpdateInput{
				InvoicingCompany: "test_company_update",
			},
			status:   "started",
			isActive: "true",
			deleted:  "false",
			mockBehavior: func(s *mock_service.MockSettlementService, settlement models.SettlementUpdateInput, id int, status, isDeleted, isActive string) {
				s.EXPECT().UpdateSettlement(id, settlement, status, isDeleted, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't update settlement","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSettlementService(c)
			testCase.mockBehavior(service, testCase.inputSettlement, testCase.inputId, testCase.status, testCase.deleted, testCase.isActive)

			controllers := NewSettlementController(service)

			r := gin.New()
			r.PATCH("/settlement/:id", controllers.UpdateSettlement)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/settlement/%b?status=%s&active=%s&deleted=%s", testCase.inputId, testCase.status, testCase.isActive, testCase.deleted),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}
*/
