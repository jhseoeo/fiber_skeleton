package req

type CommonReq struct {
	Data any `json:"data"`
}

type PaginationReq struct {
	Page  int `json:"page"  query:"page"  validate:"min=1"`
	Limit int `json:"limit" query:"limit" validate:"min=1,max=100"`
}
