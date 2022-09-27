package test

/*
func TestGroupRepository_CreateGroup(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewGroupRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Groups
		want    string
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `groups`  WHERE (company_id = ? AND name = ?)").
					WithArgs("0", "test_name").WillReturnError(errors.New("not found"))

				mock.ExpectExec("INSERT INTO `groups` (`name`,`company_id`,`priveleges`,`users`,`updated`,`created`,`is_deleted`) VALUES (?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Groups{
				Name:       "test_name",
				CompanyId:  "0",
				Priveleges: []string{"1", "2", "3"},
			},
			want: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(tt.input)
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

func TestGroupRepository_GetGroupById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewGroupRepository(nil, gdb)

	type args struct {
		GroupId   string
		companyID string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Groups
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `name`, `company_id`, `priveleges`}).
					AddRow("1", "test_name", 1, "{1,2,3,4}")

				mock.ExpectQuery("SELECT * FROM `groups` WHERE (id = ? AND company_id = ?)").
					WithArgs("1", "1").WillReturnRows(rows)
			},

			input: args{"1", "1"},

			want: models.Groups{
				ID:         1,
				Name:       "test_name",
				CompanyId:  "1",
				Priveleges: []string{"1", "2", "3", "4"},
			},
		},
		{
			name: "Groups not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `groups` WHERE (id = ? AND company_id = ?)").
					WithArgs("1", "1").WillReturnError(errors.New("group not found"))
			},
			input:   args{"1", "1"},
			wantErr: true,
			want:    models.Groups{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetGroupById(tt.input.GroupId, tt.input.companyID)
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

func TestGroupRepository_GetAllGroups(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewGroupRepository(nil, gdb)

	tests := []struct {
		name    string
		mock    func()
		want    []models.Groups
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `name`, `company_id`, `priveleges`}).
					AddRow("1", "test_name", "1", "{1,2,3}")

				mock.ExpectQuery("SELECT * FROM `groups` WHERE (id > (?))").
					WithArgs(0).WillReturnRows(rows)

			},
			want: []models.Groups{{
				ID:         1,
				Name:       "test_name",
				CompanyId:  "1",
				Priveleges: []string{"1", "2", "3"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllGroups()
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

func TestGroupRepository_DeleteGroup(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewGroupRepository(nil, gdb)

	type args struct {
		groupId string
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
				mock.ExpectExec("UPDATE `groups` SET `is_deleted` = ? WHERE (id = ?)").
					WithArgs(true, "1").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input:   args{"1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.groupId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGroupRepository_UpdateGroup(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewGroupRepository(nil, gdb)

	type args struct {
		groupId string
		group   models.GroupUpdateInput
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

				mock.ExpectExec("UPDATE `groups` SET `name` = ?, `priveleges` = ?  WHERE (id = ?)").
					WithArgs("test_name", sqlmock.AnyArg(), "1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input:   args{"1", models.GroupUpdateInput{Name: "test_name", Priveleges: []string{"1", "2", "3"}}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			_, err := r.Update(tt.input.groupId, tt.input.group)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGroupRepository_GetAllGroupsByCompanyId(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewGroupRepository(nil, gdb)

	type args struct {
		companyId string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Groups
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`id`, `name`, `company_id`, `priveleges`}).
					AddRow("1", "test_name", 1, "{1,2,3}")

				mock.ExpectQuery("SELECT * FROM `groups` WHERE (company_id = ?)").
					WithArgs("1").WillReturnRows(rows)

			},

			input: args{"1"},
			want: []models.Groups{{
				ID:         1,
				Name:       "test_name",
				CompanyId:  "1",
				Priveleges: []string{"1", "2", "3"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllGroupsByCompanyId(tt.input.companyId)
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
