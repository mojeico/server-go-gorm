package test

/*
func TestCustomerRepository_CreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewCustomerRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Customer
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `customers` (`created_at`,`updated_at`,`deleted_at`,`company_name`,`address`,`city`,`zip`,`country`,`phone_number`,`fax_number`,`legal_name`,`mc_number`,`dot_number`,`billing_address`,`billing_city`,`billing_zip`,`billing_country`,`billing_method`,`billing_type`,`billing_email`,`billing_credit_limit`,`billing_balance`,`billing_total_balance`,`status`,`is_deleted`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Customer{

				CompanyName: "CompanyName",
				Address:     "Address",
				City:        "City",

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

			got, err := r.CreateCustomer(tt.input)
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

func TestCustomerRepository_GetCustomerById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewCustomerRepository(nil, gdb)

	type args struct {
		customerId int
		status     string
		isDeleted  string
		isActive   string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Customer
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`company_name`, `address`, `city`}).
					AddRow("CompanyName", "Address", "City")

				mock.ExpectQuery("SELECT * FROM `customers` WHERE `customers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?)) ORDER BY `customers`.`id` ASC LIMIT 1").
					WithArgs(1, "started", "false", "true").WillReturnRows(rows)

			},

			input: args{1, "started", "false", "true"},

			want: models.Customer{
				CompanyName: "CompanyName",
				Address:     "Address",
				City:        "City",
			},
		},
		{
			name: "Customer not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `customers` WHERE `customers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?)) ORDER BY `customers`.`id` ASC LIMIT 1").
					WithArgs(1, "started", "false", "true").WillReturnError(errors.New("customer not found"))

			},

			input:   args{1, "started", "false", "true"},
			wantErr: true,
			want:    models.Customer{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetCustomerById(tt.input.customerId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestCustomerRepository_GetCustomers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewCustomerRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Customer
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`company_name`, `address`, `city`}).
					AddRow("CompanyName", "Address", "City")

				mock.ExpectQuery("SELECT * FROM `customers` WHERE `customers`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

			},

			input: args{"started", "false", "true"},

			want: []models.Customer{models.Customer{
				CompanyName: "CompanyName",
				Address:     "Address",
				City:        "City",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllCustomers(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestCustomerRepository_DeleteCustomer(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewCustomerRepository(nil, gdb)

	type args struct {
		customerId int
		status     string
		isDeleted  string
		isActive   string
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
				mock.ExpectExec("UPDATE `customers` SET `is_deleted` = ?, `updated_at` = ? WHERE `customers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, "started", "true").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteCustomer(tt.input.customerId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCustomerRepository_UpdateCustomer(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewCustomerRepository(nil, gdb)

	type args struct {
		customerId int
		status     string
		isDeleted  string
		isActive   string
		customer   models.CustomerUpdateInput
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
				mock.ExpectExec("UPDATE `customers` SET `legal_name` = ?, `updated_at` = ? WHERE `customers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("test", sqlmock.AnyArg(), 1, "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input:   args{1, "started", "false", "true", models.CustomerUpdateInput{LegalName: "test"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateCustomer(tt.input.customerId, tt.input.status, tt.input.isDeleted, tt.input.isActive, tt.input.customer)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
*/
