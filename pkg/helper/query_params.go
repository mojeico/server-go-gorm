package helper

type GetAllQueryParams struct {
	OffSet    string `json:"offset"`
	Limit     string `json:"limit"`
	Status    string `json:"status"`
	IsDeleted string `json:"is_deleted"`
	IsActive  string `json:"is_active"`

	QueryField string `json:"field"`
	FieldValue string `json:"value"`
}

func NewGetAllQueryParams(offSet, limit, status, isDeleted, isActive, queryField, fieldValue string) *GetAllQueryParams {
	return &GetAllQueryParams{
		OffSet:     offSet,
		Limit:      limit,
		Status:     status,
		IsDeleted:  isDeleted,
		IsActive:   isActive,
		QueryField: queryField,
		FieldValue: fieldValue,
	}
}

type GetAllQueryParamsWithId struct {
	Id        string `json:"id"`
	OffSet    string `json:"offset"`
	Limit     string `json:"limit"`
	Status    string `json:"status"`
	IsDeleted string `json:"is_deleted"`
	IsActive  string `json:"is_active"`
}

func NewGetAllQueryParamsWithId(id, offSet, limit, status, isDeleted, isActive string) *GetAllQueryParamsWithId {
	return &GetAllQueryParamsWithId{
		Id:        id,
		OffSet:    offSet,
		Limit:     limit,
		Status:    status,
		IsDeleted: isDeleted,
		IsActive:  isActive,
	}
}

type OneQueryParams struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	IsDeleted string `json:"is_deleted"`
	IsActive  string `json:"is_active"`
}

func NewOneQueryParams(id, status, isDeleted, isActive string) *OneQueryParams {
	return &OneQueryParams{
		Id:        id,
		Status:    status,
		IsDeleted: isDeleted,
		IsActive:  isActive,
	}
}

type PaginationParams struct {
	OffSet    string `json:"offset"`
	Limit     string `json:"limit"`
	IsDeleted string `json:"is_deleted"`
}

func NewPaginationParams(isDeleted, offSet, limit string) *PaginationParams {
	return &PaginationParams{
		OffSet:    offSet,
		Limit:     limit,
		IsDeleted: isDeleted,
	}
}

type InvoicingFilter struct {
	CompanyName  string `json:"company_name"`
	DeliveryFrom string `json:"delivery_from"`
	DeliveryTo   string `json:"delivery_to"`

	IsDeleted string `json:"is_deleted"`
	OffSet    string `json:"offset"`
	Limit     string `json:"limit"`
}

func NewInvoicingFilter(companyName, deliveryFrom, deliveryTo, offSet, limit, isDeleted string) *InvoicingFilter {
	return &InvoicingFilter{
		CompanyName: companyName,

		DeliveryFrom: deliveryFrom,
		DeliveryTo:   deliveryTo,

		IsDeleted: isDeleted,
		OffSet:    offSet,
		Limit:     limit,
	}
}
