package commons

// PagingResponse response de paginacion
type PagingResponse struct {
	Pagination Pagination  `json:"pagination"`
	Data       interface{} `json:"data"`
}

// Pagination atributos de paginacion
type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total,omitempty"`
}

// PagingRequest request de paginacion
type PagingRequest struct {
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	All    bool        `json:"all"`
	Data   interface{} `json:"data"`
}

// NewPagingResponse nuevo paging response
func NewPagingResponse(limit, offset, total int, data interface{}) *PagingResponse {
	return &PagingResponse{
		Pagination: Pagination{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Data: data,
	}
}

// NewPagingRequest nuevo paging request
func NewPagingRequest(limit, offset int, all bool, data interface{}) *PagingRequest {
	return &PagingRequest{
		Limit:  limit,
		Offset: offset,
		All:    all,
		Data:   data,
	}
}
