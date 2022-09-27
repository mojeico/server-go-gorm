package test

/*
var userCreateOkBody string = `{
    "company_id": 0,
    "name": "testName",
    "email": "test@gmail.com",
    "password": "testPassword",
    "auth_status":false
}`

var userUpdateOkBody string = `{
				"name": "testName",
				"groups": [1,2]
        }`

func Test_userController_CreateUser(t *testing.T) {

	type mockBehavior func(s *mock_service.MockUserService, user models.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Created",
			inputBody: userCreateOkBody,
			inputUser: models.User{
				CompanyID:           0,
				Name:                "testName",
				Email:               "test@gmail.com",
				Password:            "testPassword",
				AuthenticatedStatus: false,
			},
			mockBehavior: func(s *mock_service.MockUserService, user models.User) {
				s.EXPECT().CreateUser(user).Return(0, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"User created succesfully, id: 0","errors":null,"data":0}`,
		},
		{
			name:      "Service or Database error Created",
			inputBody: userCreateOkBody,
			inputUser: models.User{
				CompanyID:           0,
				Name:                "testName",
				Email:               "test@gmail.com",
				Password:            "testPassword",
				AuthenticatedStatus: false,
			},
			mockBehavior: func(s *mock_service.MockUserService, user models.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"Can't create a user","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			testCase.mockBehavior(service, testCase.inputUser)

			controller := controllers.NewUserController(service)

			r := gin.New()
			r.POST("/users/", controller.CreateUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_userController_DeleteUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputUserId          int
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
			inputUserId:    1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",

			mockBehavior: func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().DeleteUser(userId, companyId, status, isDeleted, isActive).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"user was deleted","errors":null,"data":1}`,
		},
		{
			name:           "Service or Database error Deleted",
			inputBody:      1,
			inputUserId:    1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().DeleteUser(userId, companyId, status, isDeleted, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"Can't delete user","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			testCase.mockBehavior(service,
				testCase.inputUserId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewUserController(service)

			r := gin.New()
			r.DELETE("/users/:id", controller.DeleteUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputUserId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_userController_GetAllUsers(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserService, status, isDeleted, isActive string)

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
			name:           "OK get all",
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockUserService, status, isDeleted, isActive string) {
				var user = models.User{}
				user.ID = 1
				s.EXPECT().GetAllUsers(status, isDeleted, isActive).Return([]models.User{user}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"Get all users successfully","errors":null,"data":[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":0,"name":"","email":"","password":"","groups":null,"auth_status":false,"status":"","is_active":false,"is_deleted":false}]}`,
		},
		{
			name:           "Service or Database error get all",
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockUserService, status, isDeleted, isActive string) {
				var user = models.User{}
				user.ID = 1
				s.EXPECT().GetAllUsers(status, isDeleted, isActive).Return([]models.User{user}, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"Can't get all users","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			testCase.mockBehavior(service,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewUserController(service)

			r := gin.New()
			r.GET("/users/", controller.GetAllUsers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/?status=%s&deleted=%s&active=%s",
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func Test_userController_GetUserById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputUserId          int
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
			mockBehavior: func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string) {
				var user = models.User{}
				user.ID = 1
				s.EXPECT().GetUserById(userId, companyId, status, isDeleted, isActive).Return(user, nil)
			},
			inputBody:            1,
			inputUserId:          1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"User was updated successfully","errors":null,"data":{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":0,"name":"","email":"","password":"","groups":null,"auth_status":false,"status":"","is_active":false,"is_deleted":false}}`,
		},
		{
			name: "Service or Database error get by id",
			mockBehavior: func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string) {
				var user = models.User{}
				user.ID = 1
				s.EXPECT().GetUserById(userId, companyId, status, isDeleted, isActive).Return(user, errors.New("service or database error"))
			},
			inputBody:            1,
			inputUserId:          1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"Can't get user by id","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			testCase.mockBehavior(service,
				testCase.inputUserId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewUserController(service)

			r := gin.New()
			r.GET("/users/user/:id", controller.GetUserById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/user/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputUserId,
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

func Test_userController_UpdateUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string, userInput models.UpdateUserInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            models.UpdateUserInput
		inputUserId          int
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
			inputBody: userUpdateOkBody,
			inputUser: models.UpdateUserInput{
				Name:   "testName",
				Groups: pq.Int32Array{1, 2},
			},
			mockBehavior: func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string, user models.UpdateUserInput) {
				s.EXPECT().UpdateUser(userId, companyId, status, isDeleted, isActive, user).Return(nil)
			},
			inputUserId:          1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"Users was updates successfully","errors":null,"data":1}`,
		},
		{
			name:           "Service or Database error Updated",
			inputBody:      userUpdateOkBody,
			inputUserId:    1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			inputUser: models.UpdateUserInput{
				Name:   "testName",
				Groups: pq.Int32Array{1, 2},
			},
			mockBehavior: func(s *mock_service.MockUserService, userId, companyId int, status, isDeleted, isActive string, user models.UpdateUserInput) {
				s.EXPECT().UpdateUser(userId, companyId, status, isDeleted, isActive, user).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"Can't update user","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			testCase.mockBehavior(service,
				testCase.inputUserId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
				testCase.inputUser,
			)

			controller := controllers.NewUserController(service)

			r := gin.New()
			r.PATCH("/users/:id", controller.UpdateUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/users/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputUserId,
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
