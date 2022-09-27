package test

/*
func TestOrderRepository_CreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewOrderRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Order
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `orders` (`created_at`,`updated_at`,`deleted_at`,`company_id`,`order_status`,`load_number`,`pickup_number`,`delivery_number`,`seal_number`,`commodity`,`weight`,`equipment_type`,`temperature_range`,`ltl`,`total_days`,`broker_load_number`,`status`,`is_deleted`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `order_shippers` (`created_at`,`updated_at`,`deleted_at`,`order_id`,`address`,`from_date`,`to_date`,`phone`) VALUES (?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `order_consignees` (`created_at`,`updated_at`,`deleted_at`,`order_id`,`consignee_address`,`from_date`,`to_date`,`phone`) VALUES (?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `order_invoices` (`created_at`,`updated_at`,`deleted_at`,`order_id`,`invoice_number`,`invoicing_company`,`bill_to_customer`,`billing_method`,`billing_type`,`driver_id`,`truck_id`,`trailer_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `order_values` (`created_at`,`updated_at`,`deleted_at`,`order_id`,`rate`,`gross_pay`,`total`,`empty_miles`,`loaded_miles`,`total_miles`,`external_notes`) VALUES (?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Order{

				CompanyId:  1,
				LoadNumber: "123",
				SealNumber: "a123",

				Status:    "started",
				IsDeleted: false,
				IsActive:  true,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateOrder(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetOrderById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewOrderRepository(nil, gdb)

	type args struct {
		orderId   int
		companyID int
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Order
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `company_id`, `load_number`, `seal_number`, `status`, `is_active`}).AddRow(0, "1", "123", "a123", "started", true)
				mock.ExpectQuery("SELECT * FROM `orders` WHERE `orders`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").WithArgs(0, 1, "started", "false", "true").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `order_consignees`  WHERE `order_consignees`.`deleted_at` IS NULL AND ((`order_id` = ?))").WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `order_shippers`  WHERE `order_shippers`.`deleted_at` IS NULL AND ((`order_id` = ?))").WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `order_invoices`  WHERE `order_invoices`.`deleted_at` IS NULL AND ((`order_id` = ?))").WithArgs(0).WillReturnRows(rows4)

				rows5 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `order_values`  WHERE `order_values`.`deleted_at` IS NULL AND ((`order_id` = ?))").WithArgs(0).WillReturnRows(rows5)

			},

			input: args{0, 1, "started", "false", "true"},

			want: models.Order{
				CompanyId:  1,
				LoadNumber: "123",
				SealNumber: "a123",

				Status:    "started",
				IsDeleted: false,
				IsActive:  true,
			},
		},
		{
			name: "Order not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `orders` WHERE `orders`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, 1, "started", "false", "true").WillReturnError(errors.New("Order not found"))

			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: true,
			want:    models.Order{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetOrderById(tt.input.orderId, tt.input.companyID, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_GetAllOrders(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewOrderRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Order
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `company_id`, `load_number`, `seal_number`, `status`, `is_active`}).
					AddRow(0, "1", "123", "a123", "started", true)

				mock.ExpectQuery("SELECT * FROM `orders` WHERE `orders`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_consignees`  WHERE `order_consignees`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_shippers`  WHERE `order_shippers`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_invoices`  WHERE `order_invoices`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows4)

				rows5 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_values`  WHERE `order_values`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows5)

			},

			input: args{"started", "false", "true"},

			want: []models.Order{{
				CompanyId:  1,
				LoadNumber: "123",
				SealNumber: "a123",

				Status:    "started",
				IsDeleted: false,
				IsActive:  true,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllOrders(tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_DeleteOrder(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewOrderRepository(nil, gdb)

	type args struct {
		OrderId   int
		companyId int
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("UPDATE `orders` SET `is_deleted` = ?, `updated_at` = ? WHERE `orders`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, 1, "started", "true").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteOrder(tt.input.OrderId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOrderRepository_UpdateOrder(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewOrderRepository(nil, gdb)

	type args struct {
		orderId   int
		companyId int
		status    string
		isDeleted string
		isActive  string
		Order     models.OrderUpdateInput
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("UPDATE `orders` SET `load_number` = ?, `updated_at` = ?  WHERE `orders`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("test number", sqlmock.AnyArg(), 0, 1, "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, 1, "started", "false", "true", models.OrderUpdateInput{LoadNumber: "test number"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateOrder(tt.input.orderId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive, tt.input.Order)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestOrderRepository_GetAllOrdersByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewOrderRepository(nil, gdb)

	type args struct {
		companyId int
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Order
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `company_id`, `load_number`, `seal_number`, `status`, `is_active`}).
					AddRow(0, "1", "123", "a123", "started", true)

				mock.ExpectQuery("SELECT * FROM `orders` WHERE `orders`.`deleted_at` IS NULL AND ((company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, "started", "false", "true").WillReturnRows(rows)

				rows1 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_consignees`  WHERE `order_consignees`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows1)

				rows2 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_shippers`  WHERE `order_shippers`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_invoices` WHERE `order_invoices`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `order_values`  WHERE `order_values`.`deleted_at` IS NULL AND ((`order_id` = ?))").
					WithArgs(0).WillReturnRows(rows4)

			},

			input: args{1, "started", "false", "true"},

			want: []models.Order{{
				CompanyId:  1,
				LoadNumber: "123",
				SealNumber: "a123",

				Status:    "started",
				IsDeleted: false,
				IsActive:  true,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllOrdersByCompanyId(tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
*/
