package test

/*
func TestTruckRepository_CreateTruck(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTruckRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Truck
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `trucks` (`created_at`,`updated_at`,`deleted_at`,`company_id`,`unit_number`,`make`,`truck_model`,`year`,`plate`,`state`,`plate_expiration`,`vin_number`,`status`,`is_active`,`is_deleted`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `truck_fuel_and_tolls` (`created_at`,`updated_at`,`deleted_at`,`truck_id`,`transporter`,`fuel_card`,`fuel_type`,`mp_g`) VALUES (?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `truck_driver_and_locations` (`created_at`,`updated_at`,`deleted_at`,`truck_id`,`driver_name`,`assigning_date`,`co_driver_name`,`trailer_name`,`location`,`location_date`,`coordinator_name`) VALUES (?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `truck_ownerships` (`created_at`,`updated_at`,`deleted_at`,`truck_id`,`company_name`,`owner_name`,`purchase_date`,`tax2290_due_date`) VALUES (?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `truck_insurances` (`created_at`,`updated_at`,`deleted_at`,`truck_id`,`liability_eff_date`,`liability_exp_date`) VALUES (?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Truck{
				CompanyId:  1,
				TruckModel: "test_model",
				UnitNumber: 123,

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

			got, err := r.CreateTruck(tt.input)
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

func TestTruckRepository_GetTruckById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTruckRepository(nil, gdb)

	type args struct {
		TruckId   int
		companyID int
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Truck
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`company_id`, `truck_model`, `unit_number`}).AddRow(1, "test_model", "123")
				mock.ExpectQuery("SELECT * FROM `trucks` WHERE `trucks`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(0, 1, "started", "false", "true").WillReturnRows(rows)

				rows1 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_fuel_and_tolls`  WHERE `truck_fuel_and_tolls`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows1)

				rows2 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_driver_and_locations`  WHERE `truck_driver_and_locations`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_ownerships`  WHERE `truck_ownerships`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_insurances`  WHERE `truck_insurances`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows4)

			},

			input: args{0, 1, "started", "false", "true"},

			want: models.Truck{
				CompanyId:  1,
				TruckModel: "test_model",
				UnitNumber: 123,

				IsDeleted: false,
				IsActive:  false,
			},
		},
		{
			name: "Truck not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `trucks` WHERE `trucks`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, 1, "started", "false", "true").WillReturnError(errors.New("Truck not found"))

			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: true,
			want:    models.Truck{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetTruckById(tt.input.TruckId, tt.input.companyID, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestTruckRepository_GetAllTrucks(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTruckRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Truck
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`company_id`, `truck_model`, `unit_number`}).AddRow(1, "test_model", "123")

				mock.ExpectQuery("SELECT * FROM `trucks` WHERE `trucks`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_fuel_and_tolls` WHERE `truck_fuel_and_tolls`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_driver_and_locations`  WHERE `truck_driver_and_locations`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_ownerships`  WHERE `truck_ownerships`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows4)

				rows5 := sqlmock.NewRows([]string{`Order_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_insurances`  WHERE `truck_insurances`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows5)

			},

			input: args{"started", "false", "true"},

			want: []models.Truck{{
				CompanyId:  1,
				TruckModel: "test_model",
				UnitNumber: 123,

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllTrucks(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestTruckRepository_DeleteTruck(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTruckRepository(nil, gdb)

	type args struct {
		TruckId   int
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
				mock.ExpectExec("UPDATE `trucks` SET `is_deleted` = ?, `updated_at` = ? WHERE `trucks`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, 1, "started", "true").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteTruck(tt.input.TruckId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTruckRepository_UpdateTruck(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTruckRepository(nil, gdb)

	type args struct {
		TruckId   int
		companyId int
		status    string
		isDeleted string
		isActive  string
		Truck     models.TruckUpdateInput
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

				mock.ExpectExec("UPDATE `trucks` SET `truck_model` = ?, `updated_at` = ? WHERE `trucks`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("test_model", sqlmock.AnyArg(), 0, 1, "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, 1, "started", "false", "true", models.TruckUpdateInput{TruckModel: "test_model"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateTruck(tt.input.TruckId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive, tt.input.Truck)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestTruckRepository_GetAllTrucksByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTruckRepository(nil, gdb)

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
		want    []models.Truck
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`company_id`, `truck_model`, `unit_number`}).AddRow(1, "test_model", "123")

				mock.ExpectQuery("SELECT * FROM `trucks` WHERE `trucks`.`deleted_at` IS NULL AND ((company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, "started", "false", "true").WillReturnRows(rows)

				rows1 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_fuel_and_tolls`  WHERE `truck_fuel_and_tolls`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows1)

				rows2 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_driver_and_locations`  WHERE `truck_driver_and_locations`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_ownerships`  WHERE `truck_ownerships`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Truck_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `truck_insurances`  WHERE `truck_insurances`.`deleted_at` IS NULL AND ((`truck_id` = ?))").
					WithArgs(0).WillReturnRows(rows4)

			},

			input: args{1, "started", "false", "true"},

			want: []models.Truck{{
				CompanyId:  1,
				TruckModel: "test_model",
				UnitNumber: 123,

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllTrucksByCompanyId(tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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
