package test

/*
var (
	orderCreateOkBody string = `{
		"load_number":"order",
		"shipper":{
			"address":"order",
			"from_date": 100,
			"to_date": 100
		},
		"values": {
			"rate": 1
		},
		"invoice": {
			"invoicing_company":"order",
			"bill_to_customer":"order",
			"billing_method":"order",
			"billing_type":"order"
		}
	}`

	orderUpdateOkBody string = `{
		"load_number":"order",
		"shipper":{
			"address":"order",
			"from_date": 100,
			"to_date": 100
		},
		"values": {
			"rate": 1
		},
		"invoice": {
			"invoicing_company":"order",
			"bill_to_customer":"order",
			"billing_method":"order",
			"billing_type":"order"
		}
	}`
)

func Test_orderController_CreateOrder(t *testing.T) {
	var orderModel = models.Order{
		LoadNumber: "order",
	}

	orderModel.Shipper = models.OrderShipper{
		Address:  "order",
		FromDate: 100,
		ToDate:   100,
	}

	orderModel.Values = models.OrderValue{
		Rate: 1,
	}

	orderModel.Invoice = models.OrderInvoice{
		InvoicingCompany: "order",
		BillToCustomer:   "order",
		BillingMethod:    "order",
		BillingType:      "order",
	}

	type mockBehavior func(s *mock_service.MockOrderService, order models.Order)

	testTable := []struct {
		name                 string
		inputBody            string
		inputOrder           models.Order
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK Created",
			inputBody:  orderCreateOkBody,
			inputOrder: orderModel,
			mockBehavior: func(s *mock_service.MockOrderService, order models.Order) {
				s.EXPECT().CreateOrder(order).Return(0, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"order created succesfully","errors":null,"data":0}`,
		},
		{
			name:       "Service or Database error Created",
			inputBody:  orderCreateOkBody,
			inputOrder: orderModel,
			mockBehavior: func(s *mock_service.MockOrderService, order models.Order) {
				s.EXPECT().CreateOrder(order).Return(0, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't create a order","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockOrderService(c)
			testCase.mockBehavior(service, testCase.inputOrder)

			controller := controllers.NewOrderController(service)

			r := gin.New()
			r.POST("/orders", controller.CreateOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/orders",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)

		})
	}
}

func Test_orderController_DeleteOrder(t *testing.T) {
	var orderModel = models.Order{
		LoadNumber: "order",
	}

	orderModel.Shipper = models.OrderShipper{
		Address:  "order",
		FromDate: 100,
		ToDate:   100,
	}

	orderModel.Values = models.OrderValue{
		Rate: 1,
	}

	orderModel.Invoice = models.OrderInvoice{
		InvoicingCompany: "order",
		BillToCustomer:   "order",
		BillingMethod:    "order",
		BillingType:      "order",
	}
	type mockBehavior func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputOrderId         int
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
			inputOrderId:   1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().DeleteOrder(orderId, companyId, status, isDeleted, isActive).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"order was deleted succesfully","errors":null,"data":1}`,
		},
		{
			name:           "Service or Database error Deleted",
			inputBody:      1,
			inputOrderId:   1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().DeleteOrder(orderId, companyId, status, isDeleted, isActive).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't delete order","errors":"service or database error","data":1}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockOrderService(c)
			testCase.mockBehavior(service,
				testCase.inputOrderId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewOrderController(service)

			r := gin.New()
			r.DELETE("/orders/:id", controller.DeleteOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/orders/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputOrderId,
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

func Test_orderController_GetAllOrdersByCompanyId(t *testing.T) {
	var orderModel = models.Order{
		LoadNumber: "order",
	}

	orderModel.ID = 1

	orderModel.Shipper = models.OrderShipper{
		Address:  "order",
		FromDate: 100,
		ToDate:   100,
	}

	orderModel.Values = models.OrderValue{
		Rate: 1,
	}

	orderModel.Invoice = models.OrderInvoice{
		InvoicingCompany: "order",
		BillToCustomer:   "order",
		BillingMethod:    "order",
		BillingType:      "order",
	}
	type mockBehavior func(s *mock_service.MockOrderService, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		inputCompanyId       int
		inputStatus          string
		inputIsDeleted       string
		inputIsActive        string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "OK get all",
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockOrderService, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().GetAllOrdersByCompanyId(companyId, status, isDeleted, isActive).Return([]models.Order{orderModel}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"orders was gotten","errors":null,"data":[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":0,"order_status":null,"load_number":"order","pickup_number":"","delivery_number":"","seal_number":"","commodity":"","weight":0,"equipment_type":"","temperature_range":null,"ltl":"","total_days":0,"broker_load_number":"","status":"","is_deleted":false,"is_active":false,"shipper":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"address":"order","from_date":100,"to_date":100,"phone":""},"consignee":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"address":"","from_date":0,"to_date":0,"phone":""},"invoice":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"invoice_number":"","invoicing_company":"order","bill_to_customer":"order","billing_method":"order","billing_type":"order","driver_id":0,"truck_id":0,"triler_id":0},"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"rate":1,"gross_pay":0,"total":0,"empty_miles":0,"loaded_miles":0,"total_miles":0,"external_notes":""}}]}`,
		},
		{
			name:           "Service or Database error get all",
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			mockBehavior: func(s *mock_service.MockOrderService, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().GetAllOrdersByCompanyId(companyId, status, isDeleted, isActive).Return([]models.Order{orderModel}, errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't get all orders","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockOrderService(c)
			testCase.mockBehavior(service,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewOrderController(service)

			r := gin.New()
			r.GET("/orders/:company_id", controller.GetAllOrdersByCompanyId)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/%d?status=%s&deleted=%s&active=%s",
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

func Test_orderController_GetOrderByID(t *testing.T) {
	var orderModel = models.Order{
		LoadNumber: "order",
	}

	orderModel.ID = 1

	orderModel.Shipper = models.OrderShipper{
		Address:  "order",
		FromDate: 100,
		ToDate:   100,
	}

	orderModel.Values = models.OrderValue{
		Rate: 1,
	}

	orderModel.Invoice = models.OrderInvoice{
		InvoicingCompany: "order",
		BillToCustomer:   "order",
		BillingMethod:    "order",
		BillingType:      "order",
	}
	type mockBehavior func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string)

	testTable := []struct {
		name                 string
		inputBody            int
		inputOrderId         int
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
			mockBehavior: func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().GetOrderById(orderId, companyId, status, isDeleted, isActive).Return(orderModel, nil)
			},
			inputBody:            1,
			inputOrderId:         1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"order by id was gotten","errors":null,"data":{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_id":0,"order_status":null,"load_number":"order","pickup_number":"","delivery_number":"","seal_number":"","commodity":"","weight":0,"equipment_type":"","temperature_range":null,"ltl":"","total_days":0,"broker_load_number":"","status":"","is_deleted":false,"is_active":false,"shipper":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"address":"order","from_date":100,"to_date":100,"phone":""},"consignee":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"address":"","from_date":0,"to_date":0,"phone":""},"invoice":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"invoice_number":"","invoicing_company":"order","bill_to_customer":"order","billing_method":"order","billing_type":"order","driver_id":0,"truck_id":0,"triler_id":0},"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"order_id":0,"rate":1,"gross_pay":0,"total":0,"empty_miles":0,"loaded_miles":0,"total_miles":0,"external_notes":""}}}`,
		},
		{
			name: "Service or Database error get by id",
			mockBehavior: func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string) {
				s.EXPECT().GetOrderById(orderId, companyId, status, isDeleted, isActive).Return(orderModel, errors.New("service or database error"))
			},
			inputBody:            1,
			inputOrderId:         1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't get order by id","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockOrderService(c)
			testCase.mockBehavior(service,
				testCase.inputOrderId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
			)

			controller := controllers.NewOrderController(service)

			r := gin.New()
			r.GET("/orders/order/:id", controller.GetOrderById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/orders/order/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputOrderId,
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

func Test_orderController_UpdateOrder(t *testing.T) {
	var orderModel = models.OrderUpdateInput{
		LoadNumber: "order",
	}

	orderModel.Shipper = models.OrderShipper{
		Address:  "order",
		FromDate: 100,
		ToDate:   100,
	}

	orderModel.Values = models.OrderValue{
		Rate: 1,
	}

	orderModel.Invoice = models.OrderInvoice{
		InvoicingCompany: "order",
		BillToCustomer:   "order",
		BillingMethod:    "order",
		BillingType:      "order",
	}

	type mockBehavior func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string, order models.OrderUpdateInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputOrder           models.OrderUpdateInput
		inputOrderId         int
		inputCompanyId       int
		inputStatus          string
		inputIsDeleted       string
		inputIsActive        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK Updated",
			inputBody:  orderUpdateOkBody,
			inputOrder: orderModel,
			mockBehavior: func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string, order models.OrderUpdateInput) {
				s.EXPECT().UpdateOrder(orderId, companyId, status, isDeleted, isActive, order).Return(nil)
			},
			inputOrderId:         1,
			inputCompanyId:       1,
			inputStatus:          "started",
			inputIsDeleted:       "false",
			inputIsActive:        "true",
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":true,"message":"order was updated succesfully","errors":null,"data":1}`,
		},
		{
			name:           "Service or Database error Updated",
			inputBody:      orderUpdateOkBody,
			inputOrderId:   1,
			inputCompanyId: 1,
			inputStatus:    "started",
			inputIsDeleted: "false",
			inputIsActive:  "true",
			inputOrder:     orderModel,
			mockBehavior: func(s *mock_service.MockOrderService, orderId, companyId int, status, isDeleted, isActive string, order models.OrderUpdateInput) {
				s.EXPECT().UpdateOrder(orderId, companyId, status, isDeleted, isActive, order).Return(errors.New("service or database error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"status":false,"message":"can't update order","errors":"service or database error","data":{}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockOrderService(c)
			testCase.mockBehavior(service,
				testCase.inputOrderId,
				testCase.inputCompanyId,
				testCase.inputStatus,
				testCase.inputIsDeleted,
				testCase.inputIsActive,
				testCase.inputOrder,
			)

			controller := controllers.NewOrderController(service)

			r := gin.New()
			r.PATCH("/orders/:id", controller.UpdateOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/orders/%d?company_id=%d&status=%s&deleted=%s&active=%s",
				testCase.inputOrderId,
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
