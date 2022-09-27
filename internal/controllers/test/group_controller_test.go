package test

/*var groupCreateOkBody string = `{
    "name": "groupName",
    "company_id": "1",
    "priveleges": ["adm","user"]
}`

var groupUpdateOkBody string = `{
				 "name": "newGroupName",
    "users": [1,2],
    "priveleges": ["new_adm","new_user"]

        }`
*/
/*
func Test_groupController_CreateGroup(t *testing.T) {

	type mockBehavior func(s *mock_service.MockGroupService, group models.Groups)

	testTable := []struct {
		name                 string
		inputBody            string
		inputGroup           models.Groups
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Created",
			inputBody: groupCreateOkBody,
			inputGroup: models.Groups{
				Name:       "groupName",
				CompanyId:  "1",
				Priveleges: pq.StringArray{"adm", "user"},
			},
			mockBehavior: func(s *mock_service.MockGroupService, group models.Groups) {
				s.EXPECT().Create(group).Return("1", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"1"`,
		},
		{
			name:      "Service or Database error Created",
			inputBody: groupCreateOkBody,
			inputGroup: models.Groups{
				Name:       "groupName",
				CompanyId:  "1",
				Priveleges: pq.StringArray{"adm", "user"},
			},
			mockBehavior: func(s *mock_service.MockGroupService, group models.Groups) {
				s.EXPECT().Create(group).Return("0", errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `"service or database error"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockGroupService(c)
			testCase.mockBehavior(service, testCase.inputGroup)

			controllers := NewGroupController(service)

			r := gin.New()
			r.POST("/groups", controllers.CreateGroup)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/groups",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_groupController_DeleteGroup(t *testing.T) {
	type mockBehavior func(s *mock_service.MockGroupService, groupId string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputGroupId         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "OK Deleted",
			inputBody:    1,
			inputGroupId: "1",

			mockBehavior: func(s *mock_service.MockGroupService, groupId string) {
				s.EXPECT().Delete(groupId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `"1"`,
		},
		{
			name:         "Service or Database error Deleted",
			inputBody:    1,
			inputGroupId: "1",
			mockBehavior: func(s *mock_service.MockGroupService, groupId string) {
				s.EXPECT().Delete(groupId).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `"service or database error"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockGroupService(c)
			testCase.mockBehavior(service, testCase.inputGroupId)

			controllers := NewGroupController(service)

			r := gin.New()
			r.DELETE("/groups/:groupId", controllers.DeleteGroup)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/groups/%s", testCase.inputGroupId),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}
func Test_groupController_GetAllGroupsByCompanyId(t *testing.T) {
	type mockBehavior func(s *mock_service.MockGroupService, companyId string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		inputCompanyId       string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "OK get all",
			inputCompanyId: "1",
			mockBehavior: func(s *mock_service.MockGroupService, companyId string) {
				s.EXPECT().GetAllGroupsByCompanyId(companyId).Return([]models.Groups{{ID: 1}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"name":"","company_id":"","priveleges":null,"users":null,"Updated":0,"Created":0,"isDeleted":false}]`,
		},
		{
			name:           "Service or Database error get all",
			inputCompanyId: "1",
			mockBehavior: func(s *mock_service.MockGroupService, companyId string) {
				s.EXPECT().GetAllGroupsByCompanyId(companyId).Return([]models.Groups{{ID: 1}}, errors.New("service or database error"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `"service or database error"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockGroupService(c)
			testCase.mockBehavior(service, testCase.inputCompanyId)

			controllers := NewGroupController(service)

			r := gin.New()
			r.GET("/groups/:companyId", controllers.GetAllGroupsByCompanyId)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/groups/%s", testCase.inputCompanyId),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func Test_groupController_GetGroupByID(t *testing.T) {
	type mockBehavior func(s *mock_service.MockGroupService, groupId, companyId string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputGroupId         string
		inputCompanyId       string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK get by id",
			mockBehavior: func(s *mock_service.MockGroupService, groupId, companyId string) {
				s.EXPECT().GetGroupById(groupId, companyId).Return(models.Groups{ID: 1}, nil)
			},
			inputBody:            1,
			inputGroupId:         "1",
			inputCompanyId:       "1",
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"name":"","company_id":"","priveleges":null,"users":null,"Updated":0,"Created":0,"isDeleted":false}`,
		},
		{
			name: "Service or Database error get by id",
			mockBehavior: func(s *mock_service.MockGroupService, groupId, companyId string) {
				s.EXPECT().GetGroupById(groupId, companyId).Return(models.Groups{ID: 1}, errors.New("service or database error"))
			},
			inputBody:            1,
			inputGroupId:         "1",
			inputCompanyId:       "1",
			expectedStatusCode:   400,
			expectedResponseBody: `"Can not find user by id"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockGroupService(c)
			testCase.mockBehavior(service, testCase.inputGroupId, testCase.inputCompanyId)

			controllers := NewGroupController(service)

			r := gin.New()
			r.GET("/groups/:groupId", controllers.GetGroupByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/groups/%s?company_id=%s", testCase.inputGroupId, testCase.inputCompanyId),
				bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_groupController_UpdateGroup(t *testing.T) {
	type mockBehavior func(s *mock_service.MockGroupService, groupInput models.GroupUpdateInput, groupId string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputGroup           models.GroupUpdateInput
		inputGroupId         string
		inputCompanyId       int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK Updated",
			inputBody: groupUpdateOkBody,
			inputGroup: models.GroupUpdateInput{
				Name:       "newGroupName",
				Users:      pq.Int32Array{1, 2},
				Priveleges: pq.StringArray{"new_adm", "new_user"},
			},
			mockBehavior: func(s *mock_service.MockGroupService, group models.GroupUpdateInput, groupId string) {
				s.EXPECT().Update(groupId, group).Return(models.Groups{ID: 1}, nil)
			},
			inputGroupId:         "1",
			inputCompanyId:       1,
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"name":"","company_id":"","priveleges":null,"users":null,"Updated":0,"Created":0,"isDeleted":false}`,
		},
		{
			name:           "Service or Database error Updated",
			inputBody:      groupUpdateOkBody,
			inputGroupId:   "1",
			inputCompanyId: 1,
			inputGroup: models.GroupUpdateInput{
				Name:       "newGroupName",
				Users:      pq.Int32Array{1, 2},
				Priveleges: pq.StringArray{"new_adm", "new_user"},
			},
			mockBehavior: func(s *mock_service.MockGroupService, group models.GroupUpdateInput, groupId string) {
				s.EXPECT().Update(groupId, group).Return(models.Groups{ID: 1}, errors.New("service or database error"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: `"service or database error"`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockGroupService(c)
			testCase.mockBehavior(service, testCase.inputGroup, testCase.inputGroupId)

			controllers := NewGroupController(service)

			r := gin.New()
			r.PATCH("/Groups/:groupId", controllers.UpdateGroup)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/Groups/%s?company_id=%b", testCase.inputGroupId, testCase.inputCompanyId),
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}*/
