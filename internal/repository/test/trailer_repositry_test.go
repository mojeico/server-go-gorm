package test

/*
func TestTrailerRepository_CreateTrailer(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTrailerRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Trailer
		want    models.Trailer
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `trailers` (`id`,`driver_id`,`name`,`type`,`unit_number`,`make`,`year`,`plate_state`,`expiration`,`vinn_umber`,`weight`,`status`,`is_active`,`is_deleted`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Trailer{

				Name:       "test_name",
				UnitNumber: "test_number",
				PlateState: "state",

				Status:    "started",
				IsDeleted: false,
				IsActive:  true,
			},
			want: models.Trailer{

				Name:       "test_name",
				UnitNumber: "test_number",
				PlateState: "state",

				Status:    "started",
				IsDeleted: false,
				IsActive:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.TrailerCreate(tt.input)
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

func TestTrailerRepository_GetTrailerById(t *testing.T) {

}

func TestTrailerRepository_GetAllTrailers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTrailerRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Trailer
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`name`, `unit_number`, `plate_state`}).AddRow("test_name", "test_number", "state")

				mock.ExpectQuery("SELECT * FROM `trailers` WHERE (status = ? AND is_deleted = ? AND is_active = ?)").
					WithArgs("started", "false", "true").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_checks`  WHERE `trailer_checks`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_inside_cargos`  WHERE `trailer_inside_cargos`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_locations`  WHERE `trailer_locations`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows4)

				rows5 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_owners`  WHERE `trailer_owners`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows5)

				rows6 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_what_evers`  WHERE `trailer_what_evers`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows6)

			},

			input: args{"started", "false", "true"},

			want: []models.Trailer{{
				Name:       "test_name",
				UnitNumber: "test_number",
				PlateState: "state",

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllTrailers(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestTrailerRepository_DeleteTrailer(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTrailerRepository(nil, gdb)

	type args struct {
		TrailerId int
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

			err := r.TrailerDelete(tt.input.TrailerId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTrailerRepository_UpdateTrailer(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTrailerRepository(nil, gdb)

	type args struct {
		TrailerId int
		companyId int
		status    string
		isDeleted string
		isActive  string
		Trailer   models.TrailerInputUpdate
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

				mock.ExpectExec("UPDATE `trailers` SET `name` = ? WHERE (id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?)").
					WithArgs("update_name", 0, 1, "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, 1, "started", "false", "true", models.TrailerInputUpdate{Name: "update_name"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.TrailerUpdate(tt.input.TrailerId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive, tt.input.Trailer)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestTrailerRepository_GetAllTrailersByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewTrailerRepository(nil, gdb)

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
		want    []models.Trailer
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`name`, `unit_number`, `plate_state`}).AddRow("test_name", "test_number", "state")

				mock.ExpectQuery("SELECT * FROM `trailers` WHERE (company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?)").
					WithArgs(1, "", "false", "false").WillReturnRows(rows)

				rows2 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_checks`  WHERE `trailer_checks`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows2)

				rows3 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_inside_cargos`  WHERE `trailer_inside_cargos`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows3)

				rows4 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_locations`  WHERE `trailer_locations`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows4)

				rows5 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_owners`  WHERE `trailer_owners`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows5)

				rows6 := sqlmock.NewRows([]string{`Order_id`}).AddRow("0")
				mock.ExpectQuery("SELECT * FROM `trailer_what_evers`  WHERE `trailer_what_evers`.`deleted_at` IS NULL AND ((`trailer_id` = ?))").WithArgs(0).WillReturnRows(rows6)

			},

			input: args{1, "", "false", "false"},

			want: []models.Trailer{{
				Name:       "test_name",
				UnitNumber: "test_number",
				PlateState: "state",

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllTrailerByCompanyId(tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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
