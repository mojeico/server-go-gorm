package test

/*
func TestDriverRepository_CreateDriver(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewDriverRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Driver
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `drivers` (`created_at`,`updated_at`,`deleted_at`,`first_name`,`last_name`,`address`,`city`,`state`,`zip`,`country`,`phone`,`email`,`gender`,`birth_day`,`status`,`is_active`,`is_deleted`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `checks` (`created_at`,`updated_at`,`deleted_at`,`driver_id`,`active`,`freeze_payable`,`gross_pay`,`eligible_for_rehire`,`checks`) VALUES (?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO `driver_attaches` (`created_at`,`updated_at`,`deleted_at`,`driver_id`,`company_id`,`driver_group_id`,`sattlement_id`,`employment_id`,`safety_id`) VALUES (?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Driver{

				FirstName: "CompanyName",
				LastName:  "Address",
				City:      "City",

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

			got, err := r.CreateDriver(tt.input)
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

func TestDriverRepository_GetDriverById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewDriverRepository(nil, gdb)

	type args struct {
		driverId  int
		companyID int
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Driver
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `first_name`, `last_name`, `city`}).
					AddRow(0, "FirstName", "LastName", "City")

				mock.ExpectQuery("SELECT * FROM `drivers` WHERE `drivers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(0, "started", "false", "true").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `driver_attaches`  WHERE `driver_attaches`.`deleted_at` IS NULL AND ((`driver_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `checks`  WHERE `checks`.`deleted_at` IS NULL AND ((`driver_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

			},

			input: args{0, 1, "started", "false", "true"},

			want: models.Driver{
				FirstName: "FirstName",
				LastName:  "LastName",
				City:      "City",
			},
		},
		{
			name: "Driver not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `drivers` WHERE `drivers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, "started", "false", "true").WillReturnError(errors.New("Driver not found"))

			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: true,
			want:    models.Driver{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetDriverById(tt.input.driverId, tt.input.companyID, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestDriverRepository_GetAllDrivers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewDriverRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Driver
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `first_name`, `last_name`, `city`}).
					AddRow(0, "FirstName", "LastName", "City")

				mock.ExpectQuery("SELECT * FROM `drivers` WHERE `drivers`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `driver_attaches`  WHERE `driver_attaches`.`deleted_at` IS NULL AND ((`driver_id` = ?))").
					WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `checks`  WHERE `checks`.`deleted_at` IS NULL AND ((`driver_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

			},

			input: args{"started", "false", "true"},

			want: []models.Driver{models.Driver{
				FirstName: "FirstName",
				LastName:  "LastName",
				City:      "City",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllDrivers(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestDriverRepository_DeleteDriver(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewDriverRepository(nil, gdb)

	type args struct {
		driverId  int
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
				mock.ExpectExec("UPDATE `drivers` SET `is_deleted` = ?, `updated_at` = ? WHERE `drivers`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, "started", "true").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteDriver(tt.input.driverId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDriverRepository_UpdateDriver(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewDriverRepository(nil, gdb)

	type args struct {
		driverId  int
		companyId int
		status    string
		isDeleted string
		isActive  string
		driver    models.DriverUpdateInput
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

				rows := sqlmock.NewRows([]string{`id`, `first_name`, `last_name`, `city`}).
					AddRow(0, "FirstName", "LastName", "City")

				mock.ExpectQuery("SELECT * FROM `drivers` WHERE `drivers`.`deleted_at` IS NULL AND ((id = ? AND is_deleted = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(0, "started", "false", "true").WillReturnRows(rows)

				mock.ExpectExec("UPDATE `drivers` SET `first_name` = ?, `updated_at` = ? WHERE `drivers`.`deleted_at` IS NULL AND ((id = ? AND is_deleted = ?))").
					WithArgs("test", sqlmock.AnyArg(), 0, false).WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, 1, "started", "false", "true", models.DriverUpdateInput{FirstName: "test"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateDriver(tt.input.driverId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive, tt.input.driver)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestDriverRepository_GetAllDriversByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewDriverRepository(nil, gdb)

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
		want    []models.Driver
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows2 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `driver_attaches` WHERE `driver_attaches`.`deleted_at` IS NULL AND ((company_id = ?))").
					WithArgs(1).WillReturnRows(rows2)

				rows := sqlmock.NewRows([]string{`id`, `first_name`, `last_name`, `city`}).
					AddRow(0, "FirstName", "LastName", "City")

				mock.ExpectQuery("SELECT * FROM `drivers` WHERE `drivers`.`deleted_at` IS NULL AND ((id = ?))").
					WithArgs(0).WillReturnRows(rows)

				rows3 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `driver_attaches`  WHERE `driver_attaches`.`deleted_at` IS NULL AND ((`driver_id` = ?))").
					WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`driver_id`}).
					AddRow("0")

				mock.ExpectQuery("SELECT * FROM `checks`  WHERE `checks`.`deleted_at` IS NULL AND ((`driver_id` = ?))").
					WithArgs(0).WillReturnRows(rows4)

			},

			input: args{1, "started", "false", "true"},

			want: []models.Driver{models.Driver{
				FirstName: "FirstName",
				LastName:  "LastName",
				City:      "City",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllDriversByCompanyId(tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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
