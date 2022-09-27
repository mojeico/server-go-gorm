package test

/*
func TestSettlementRepository_CreateSettlement(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSettlementRepository(gdb)

	tests := []struct {
		name    string
		mock    func()
		input   models.Settlement
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("INSERT INTO `settlements` (`created_at`,`updated_at`,`deleted_at`,`settlement_date`,`invoicing_company`,`driver`,`total_miles`,`empty_miles`,`loaded_miles`,`date_submitted`,`deductions`,`reimbursement`,`earning`,`total`,`status`,`is_deleted`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input: models.Settlement{

				InvoicingCompany: "test_company",
				Driver:           123321,
				Total:            111.111,

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

			got, err := r.CreateSettlement(tt.input)
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

func TestSettlementRepository_GetSettlementById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSettlementRepository(gdb)

	type args struct {
		SettlementId int

		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.Settlement
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`invoicing_company`, `driver`, `total`}).
					AddRow("test_company", 123321, 111.111)

				mock.ExpectQuery("SELECT * FROM `settlements` WHERE `settlements`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?)) ORDER BY `settlements`.`id` ASC LIMIT 1").
					WithArgs(1, "started", "false", "true").WillReturnRows(rows)
			},

			input: args{1, "started", "false", "true"},

			want: models.Settlement{
				InvoicingCompany: "test_company",
				Driver:           123321,
				Total:            111.111,

				IsDeleted: false,
				IsActive:  false,
			},
		},
		{
			name: "Groups not found",
			mock: func() {

				mock.ExpectQuery("SELECT * FROM `settlements` WHERE `settlements`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?)) ORDER BY `settlements`.`id` ASC LIMIT 1").
					WithArgs(1, "started", "false", "true").WillReturnError(errors.New("group not found"))
			},
			input:   args{1, "started", "false", "true"},
			wantErr: true,
			want:    models.Settlement{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSettlementById(tt.input.SettlementId, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestSettlementRepository_GetAllSettlements(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSettlementRepository(gdb)

	type args struct {
		status    string
		isDeleted string
		isActive  string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Settlement
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				rows := sqlmock.NewRows([]string{`invoicing_company`, `driver`, `total`}).
					AddRow("test_company", 123321, 111.111)

				mock.ExpectQuery("SELECT * FROM `settlements` WHERE `settlements`.`deleted_at` IS NULL AND ((status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs("started", "false", "true").WillReturnRows(rows)

			},

			input: args{"started", "false", "true"},

			want: []models.Settlement{{

				InvoicingCompany: "test_company",
				Driver:           123321,
				Total:            111.111,

				IsDeleted: false,
				IsActive:  false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAllSettlement(tt.input.status, tt.input.isDeleted, tt.input.isActive)
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

func TestSettlementRepository_DeleteSettlement(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSettlementRepository(gdb)

	type args struct {
		SettlementId int
		status       string
		isDeleted    string
		isActive     string
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
				mock.ExpectExec("UPDATE `settlements` SET `is_deleted` = ?, `updated_at` = ? WHERE `settlements`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_active = ?))").
					WithArgs(true, sqlmock.AnyArg(), 1, "started", "true").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:   args{1, "started", "false", "true"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteSettlement(tt.input.SettlementId, tt.input.status, tt.input.isActive)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSettlementRepository_UpdateSettlement(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	gdb, err := gorm.Open("mysql", db)

	r := repository.NewSettlementRepository(gdb)

	type args struct {
		SettlementId int
		companyId    int
		status       string
		isDeleted    string
		isActive     string
		Settlement   models.SettlementUpdateInput
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

				mock.ExpectExec("UPDATE `settlements` SET `driver` = ?, `updated_at` = ? WHERE `settlements`.`deleted_at` IS NULL AND ((id = ? AND status = ? AND is_deleted = ? AND is_active = ?))").
					WithArgs(123, sqlmock.AnyArg(), 0, "started", "false", "true").WillReturnResult(sqlmock.NewResult(1, 1))

			},
			input:   args{0, 1, "started", "false", "true", models.SettlementUpdateInput{Driver: 123}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.UpdateSettlement(tt.input.SettlementId, tt.input.Settlement, tt.input.status, tt.input.isDeleted, tt.input.isActive)
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
