package test

/*
func TestSafetyRepository_CreateSafety(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSafetyRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Safety
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `safeties` (`created_at`,`updated_at`,`deleted_at`,`company_id`,`uploading_date`,`file_type`,`file_name`,`expiration_date`,`comments`,`status`,`is_deleted`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Safety{

				FileType: "test_file_type",
				FileName: "test_file_name",

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

			got, err := r.CreateSafety(tt.input)
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

func TestSafetyRepository_GetSafetyById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSafetyRepository(nil, gdb)

	type args struct {
		SafetyId  int
		companyID string
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Safety
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `file_type`, `file_name`}).
					AddRow(0, "test_file_type", "test_file_name")

				mock.ExpectQuery("SELECT * FROM `safeties` WHERE `safeties`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?)) ORDER BY `safeties`.`id` ASC LIMIT 1").
					WithArgs(0, "1", "started", "false", "true").WillReturnRows(rows)

			},

			input: args{0, "1", "started", "false", "true"},

			want: models.Safety{

				FileType: "test_file_type",
				FileName: "test_file_name",

				IsDeleted: false,
				IsActive:  false,
			},
		},
		{
			name: "Safety not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `safeties` WHERE `safeties`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?)) ORDER BY `safeties`.`id` ASC LIMIT 1").
					WithArgs(0, "1", "started", "false", "true").WillReturnError(errors.New("Driver not found"))

			},
			input:   args{0, "1", "started", "false", "true"},
			wantErr: true,
			want:    models.Safety{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSafetyById(tt.input.SafetyId, tt.input.companyID, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestSafetyRepository_GetAllSafetys(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSafetyRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Safety
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `file_type`, `file_name`}).
					AddRow(0, "test_file_type", "test_file_name")

				mock.ExpectQuery("SELECT * FROM `safeties` WHERE `safeties`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

			},

			input: args{"started", "false", "true"},

			want: []models.Safety{models.Safety{

				FileType: "test_file_type",
				FileName: "test_file_name",

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllSafeties(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestSafetyRepository_DeleteSafety(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSafetyRepository(nil, gdb)

	type args struct {
		SafetyId  int
		companyId string
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
				mock.ExpectExec("UPDATE `safeties` SET `is_deleted` = ?, `updated_at` = ? WHERE `safeties`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, "1", "started", "true").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, "1", "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteSafety(tt.input.SafetyId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSafetyRepository_UpdateSafety(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSafetyRepository(nil, gdb)

	type args struct {
		SafetyId  int
		companyId string
		status    string
		isDeleted string
		isActive  string
		Safety    models.SafetyInputUpdate
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

				mock.ExpectExec("UPDATE `safeties` SET `file_name` = ?, `updated_at` = ?  WHERE `safeties`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("test_new_file_name", sqlmock.AnyArg(), 0, "1", "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, "1", "started", "false", "true", models.SafetyInputUpdate{FileName: "test_new_file_name"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateSafety(tt.input.SafetyId, tt.input.Safety, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestSafetyRepository_GetAllSafetysByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSafetyRepository(nil, gdb)

	type args struct {
		companyId string
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Safety
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `file_type`, `file_name`}).
					AddRow("0", "test_file_type", "test_file_name")

				mock.ExpectQuery("SELECT * FROM `safeties` WHERE `safeties`.`deleted_at` IS NULL AND ((company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("0", "started", "false", "true").WillReturnRows(rows)

			},

			input: args{"0", "started", "false", "true"},

			want: []models.Safety{models.Safety{

				FileType: "test_file_type",
				FileName: "test_file_name",

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllSafetiesByCompanyId(tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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
