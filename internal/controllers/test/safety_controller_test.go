package test

/*
var safetyCreateOkBody string = `{
        "company_id" : 1,
		"uploading_date": 123123,
		"file_type" : "FileType",
		"file_name" : "FileName",
		"expiration_date" : 123123
	 }`

var safetyUpdateOkBody string = `{
           "file_name" : "new test name"
        }`

func TestSafetyController_CreateSafety(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSafetyService, safety models.Safety)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSafety          models.Safety
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Created",
			inputBody: safetyCreateOkBody,
			inputSafety: models.Safety{
				CompanyID:      1,
				UploadingDate:  123123,
				FileType:       "FileType",
				FileName:       "FileName",
				ExpirationDate: 123123,
			},
			mockBehavior: func(s *mock_service.MockSafetyService, safety models.Safety) {
				s.EXPECT().CreateSafety(safety).Return(0, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"status":true,"message":"Safety created successfully, id: 0","errors":null,"data":0}`,
		},
		{
			name:      "Service or Database error created",
			inputBody: safetyCreateOkBody,
			inputSafety: models.Safety{
				CompanyID:      1,
				UploadingDate:  123123,
				FileType:       "FileType",
				FileName:       "FileName",
				ExpirationDate: 123123,
			},
			mockBehavior: func(s *mock_service.MockSafetyService, safety models.Safety) {
				s.EXPECT().CreateSafety(safety).Return(0, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't create safety ","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSafetyService(c)
			testCase.mockBehavior(service, testCase.inputSafety)

			controller := controllers.NewSafetyController(service)

			r := gin.New()
			r.POST("/safety", controller.CreateSafety)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/safety",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func TestSafetyController_DeleteSafety(t *testing.T) {

	type mockBehavior func(s *mock_service.MockSafetyService, id int, companyId, status, isDelete, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputSafety          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
		companyId            string
	}{
		{
			name:        "OK Deleted",
			inputBody:   1,
			inputSafety: 1,
			companyId:   "1",
			mockBehavior: func(s *mock_service.MockSafetyService, id int, companyId, status, isDelete, isActive string) {
				s.EXPECT().DeleteSafety(id, companyId, status, isDelete, isActive).Return(nil)
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"Safety deleted successfully, id: 1","errors":null,"data":1}`,
		},
		{
			name:        "Service or Database error Deleted",
			inputBody:   1,
			inputSafety: 1,
			companyId:   "1",
			status:      "started",
			isActive:    "true",
			deleted:     "false",

			mockBehavior: func(s *mock_service.MockSafetyService, id int, companyId, status, isDelete, isActive string) {
				s.EXPECT().DeleteSafety(id, companyId, status, isDelete, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can not delete safety","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSafetyService(c)
			testCase.mockBehavior(service, testCase.inputSafety, testCase.companyId, testCase.status, testCase.deleted, testCase.isActive)

			controller := controllers.NewSafetyController(service)

			r := gin.New()
			r.DELETE("/safety/:id", controller.DeleteSafety)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/safety/%b?status=%s&active=%s&deleted=%s&companyId=%s", testCase.inputBody, testCase.status, testCase.isActive, testCase.deleted, testCase.companyId),
				bytes.NewBufferString(fmt.Sprint(testCase.inputBody)))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}

}

func TestSafetyController_GetAllSafeties(t *testing.T) {

	type mockBehavior func(s *mock_service.MockSafetyService, status, isDeleted, isActive string)

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

			mockBehavior: func(s *mock_service.MockSafetyService, status, isDeleted, isActive string) {
				s.EXPECT().GetAllSafeties(status, isDeleted, isActive).Return([]models.Safety{{CompanyID: 1}}, nil)
			},
			expectedStatusCode:   200,
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedResponseBody: `{"status":true,"message":"Return all safeties","errors":null,"data":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":1,"uploading_date":0,"file_type":"","file_name":"","expiration_date":0,"comments":"","status":"","deleted":false,"is_active":false}]}`,
		},
		{
			name: "Service or Database error get all",
			mockBehavior: func(s *mock_service.MockSafetyService, status, isDeleted, isActive string) {
				s.EXPECT().GetAllSafeties(status, isDeleted, isActive).Return([]models.Safety{{CompanyID: 1}}, errors.New("service or database error"))
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't get all safeties","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSafetyService(c)
			testCase.mockBehavior(service, testCase.status, testCase.deleted, testCase.isActive)

			controller := controllers.NewSafetyController(service)

			r := gin.New()
			r.GET("/safety", controller.GetAllSafeties)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/safety?status=%s&active=%s&deleted=%s", testCase.status, testCase.isActive, testCase.deleted),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}

}

func TestSafetyController_GetAllSafetiesByCompanyId(t *testing.T) {

	type mockBehavior func(s *mock_service.MockSafetyService, companyId, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
		companyId            string
	}{
		{
			name: "OK get all",

			mockBehavior: func(s *mock_service.MockSafetyService, companyId, status, isDeleted, isActive string) {
				s.EXPECT().GetAllSafetiesByCompanyId(companyId, status, isDeleted, isActive).Return([]models.Safety{{CompanyID: 1}}, nil)
			},
			expectedStatusCode:   200,
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			companyId:            "1",
			expectedResponseBody: `{"status":true,"message":"Get all safeties by company id successfully","errors":null,"data":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":1,"uploading_date":0,"file_type":"","file_name":"","expiration_date":0,"comments":"","status":"","deleted":false,"is_active":false}]}`,
		},
		{
			name: "Service or Database error get all",
			mockBehavior: func(s *mock_service.MockSafetyService, companyId, status, isDeleted, isActive string) {
				s.EXPECT().GetAllSafetiesByCompanyId(companyId, status, isDeleted, isActive).Return([]models.Safety{{CompanyID: 1}}, errors.New("service or database error"))
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			expectedStatusCode:   500,
			companyId:            "1",
			expectedResponseBody: `{"status":false,"message":"Can't get all safeties by company id","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSafetyService(c)
			testCase.mockBehavior(service, testCase.companyId, testCase.status, testCase.deleted, testCase.isActive)

			controller := controllers.NewSafetyController(service)

			r := gin.New()
			r.GET("/safety/:companyId", controller.GetAllSafetiesByCompanyId)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/safety/%s?status=%s&active=%s&deleted=%s", testCase.companyId, testCase.status, testCase.isActive, testCase.deleted),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}

}

func TestSafetyController_GetSafetyById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSafetyService, id int, companyId, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSafety          models.SafetyInputUpdate
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
		companyId            string
	}{
		{
			name: "OK get by id",
			mockBehavior: func(s *mock_service.MockSafetyService, id int, companyId, status, isDeleted, isActive string) {
				s.EXPECT().GetSafetyById(id, companyId, status, isDeleted, isActive).Return(models.Safety{CompanyID: 1}, nil)
			},
			status:               "started",
			isActive:             "true",
			deleted:              "false",
			inputId:              1,
			expectedStatusCode:   200,
			companyId:            "1",
			expectedResponseBody: `{"status":true,"message":"Return one safety","errors":null,"data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":1,"uploading_date":0,"file_type":"","file_name":"","expiration_date":0,"comments":"","status":"","deleted":false,"is_active":false}}`,
		},
		{
			name: "Service or Database error get by id",
			mockBehavior: func(s *mock_service.MockSafetyService, id int, companyId, status, isDeleted, isActive string) {
				s.EXPECT().GetSafetyById(id, companyId, status, isDeleted, isActive).Return(models.Safety{}, errors.New("service or database error"))
			},
			status:    "started",
			isActive:  "true",
			deleted:   "false",
			inputId:   1,
			companyId: "1",

			expectedStatusCode:   400,
			expectedResponseBody: `{"status":false,"message":"can't get safety by id","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSafetyService(c)
			testCase.mockBehavior(service, testCase.inputId, testCase.companyId, testCase.status, testCase.deleted, testCase.isActive)

			controller := controllers.NewSafetyController(service)

			r := gin.New()
			r.GET("/safety/:id", controller.GetSafetyById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/safety/%b?status=%s&active=%s&deleted=%s&companyId=%s", testCase.inputId, testCase.status, testCase.isActive, testCase.deleted, testCase.companyId),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func TestSafetyController_UpdateSafety(t *testing.T) {
	type mockBehavior func(s *mock_service.MockSafetyService, safety models.SafetyInputUpdate, id int, companyId, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputSafety          models.SafetyInputUpdate
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		status               string
		isActive             string
		deleted              string
		companyId            string
	}{
		{
			name:      "OK Updated",
			inputBody: safetyUpdateOkBody,
			inputSafety: models.SafetyInputUpdate{
				FileName: "new test name",
			},
			status:   "started",
			isActive: "true",
			deleted:  "false",
			mockBehavior: func(s *mock_service.MockSafetyService, safety models.SafetyInputUpdate, id int, companyId, status, isDeleted, isActive string) {
				s.EXPECT().UpdateSafety(id, safety, companyId, status, isDeleted, isActive).Return(nil)
			},
			inputId:              1,
			expectedStatusCode:   200,
			companyId:            "1",
			expectedResponseBody: `{"status":true,"message":"Update safety successfully","errors":null,"data":{}}`,
		},
		{
			name:      "Service or Database error Updated",
			inputBody: safetyUpdateOkBody,
			inputSafety: models.SafetyInputUpdate{
				FileName: "new test name",
			},
			status:    "started",
			isActive:  "true",
			deleted:   "false",
			companyId: "1",
			mockBehavior: func(s *mock_service.MockSafetyService, safety models.SafetyInputUpdate, id int, companyId, status, isDeleted, isActive string) {
				s.EXPECT().UpdateSafety(id, safety, companyId, status, isDeleted, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"Can't update safety","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockSafetyService(c)
			testCase.mockBehavior(service, testCase.inputSafety, testCase.inputId, testCase.companyId, testCase.status, testCase.deleted, testCase.isActive)

			controller := controllers.NewSafetyController(service)

			r := gin.New()
			r.PATCH("/safety/:id", controller.UpdateSafety)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/safety/%b?status=%s&active=%s&deleted=%s&companyId=%s", testCase.inputId, testCase.status, testCase.isActive, testCase.deleted, testCase.companyId),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}
*/
