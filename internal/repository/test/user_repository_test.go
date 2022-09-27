package test

/*
func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewUserRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`company_id`,`name`,`email`,`password`,`groups`,`authenticated_status`,`status`,`is_active`,`is_deleted`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.User{
				Name:  "test_name",
				Email: "test_email",

				Status:    "started",
				IsDeleted: false,
				IsActive:  false,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
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

func TestUserRepository_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewUserRepository(nil, gdb)

	type args struct {
		UserId    int
		companyID int
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`name`, `email`}).AddRow("test_name", "test_email")
				mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(0, 1, "started", "false", "true").WillReturnRows(rows)

			},

			input: args{0, 1, "started", "false", "true"},

			want: models.User{
				Name:  "test_name",
				Email: "test_email",

				IsDeleted: false,
				IsActive:  false,
			},
		},
		{
			name: "User not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, 1, "started", "false", "true").WillReturnError(errors.New("user not found"))

			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: true,
			want:    models.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserById(tt.input.UserId, tt.input.companyID, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestUserRepository_GetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewUserRepository(nil, gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`name`, `email`}).AddRow("test_name", "test_email")

				mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

			},

			input: args{"started", "false", "true"},

			want: []models.User{{
				Name:  "test_name",
				Email: "test_email",

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllUsers(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestUserRepository_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewUserRepository(nil, gdb)

	type args struct {
		UserId    int
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
				mock.ExpectExec("UPDATE `users` SET `is_deleted` = ?, `updated_at` = ? WHERE `users`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, 1, "started", "false", "true").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, 1, "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteUser(tt.input.UserId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewUserRepository(nil, gdb)

	type args struct {
		UserId    int
		companyId int
		status    string
		isDeleted string
		isActive  string
		User      models.UpdateUserInput
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

				mock.ExpectExec("UPDATE `users` SET `name` = ?, `updated_at` = ? WHERE `users`.`deleted_at` IS NULL AND ((id = ? AND company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("test_name", sqlmock.AnyArg(), 0, 1, "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, 1, "started", "false", "true", models.UpdateUserInput{Name: "test_name"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateUser(tt.input.UserId, tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive, tt.input.User)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestUserRepository_GetAllUsersByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewUserRepository(nil, gdb)

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
		want    []models.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`name`, `email`}).AddRow("test_name", "test_email")

				mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL AND ((company_id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(1, "started", "false", "true").WillReturnRows(rows)

			},

			input: args{1, "started", "false", "true"},

			want: []models.User{{
				Name:  "test_name",
				Email: "test_email",

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllUsersByCompanyId(tt.input.companyId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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
